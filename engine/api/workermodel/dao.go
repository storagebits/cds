package workermodel

import (
	"context"
	"time"

	"github.com/go-gorp/gorp"
	"github.com/lib/pq"

	"github.com/ovh/cds/engine/api/database/gorpmapping"
	"github.com/ovh/cds/engine/api/group"
	"github.com/ovh/cds/engine/api/secret"
	"github.com/ovh/cds/sdk"
	"github.com/ovh/cds/sdk/log"
)

func getAllCommon(ctx context.Context, db gorp.SqlExecutor, q gorpmapping.Query, opts ...LoadOptionFunc) (sdk.Models, error) {
	ms := []*workerModel{}

	if err := gorpmapping.GetAll(ctx, db, q, &ms); err != nil {
		return nil, sdk.WrapError(err, "cannot get worker models")
	}

	// Check signature of data, if invalid do not return it
	verifiedModels := make([]*sdk.Model, 0, len(ms))
	for i := range ms {
		isValid, err := gorpmapping.CheckSignature(ms[i], ms[i].Signature)
		if err != nil {
			return nil, err
		}
		if !isValid {
			log.Error(ctx, "workermodel.getAll> worker model %d data corrupted", ms[i].ID)
			continue
		}
		verifiedModels = append(verifiedModels, &ms[i].Model)
	}

	if len(verifiedModels) > 0 {
		for i := range opts {
			if err := opts[i](ctx, db, verifiedModels...); err != nil {
				return nil, err
			}
		}
	}

	models := make([]sdk.Model, len(verifiedModels))
	for i := range verifiedModels {
		models[i] = *verifiedModels[i]
	}

	return models, nil
}

func getAll(ctx context.Context, db gorp.SqlExecutor, q gorpmapping.Query, opts ...LoadOptionFunc) (sdk.Models, error) {
	ms, err := getAllCommon(ctx, db, q, opts...)
	if err != nil {
		return nil, err
	}

	// Temporary hide password
	for i := range ms {
		if ms[i].ModelDocker.Password != "" {
			ms[i].ModelDocker.Password = sdk.PasswordPlaceholder
		}
	}

	return ms, nil
}

func getAllWithClearPassword(ctx context.Context, db gorp.SqlExecutor, q gorpmapping.Query, opts ...LoadOptionFunc) ([]sdk.Model, error) {
	ms, err := getAllCommon(ctx, db, q, opts...)
	if err != nil {
		return nil, err
	}

	// Temporary decrypt password
	for i := range ms {
		if ms[i].ModelDocker.Private && ms[i].ModelDocker.Password != "" {
			var err error
			ms[i].ModelDocker.Password, err = secret.DecryptValue(ms[i].ModelDocker.Password)
			if err != nil {
				return nil, sdk.WrapError(err, "cannot decrypt value for model with id %d", ms[i].ID)
			}
		}
	}

	return ms, nil
}

func getCommon(ctx context.Context, db gorp.SqlExecutor, q gorpmapping.Query, opts ...LoadOptionFunc) (*sdk.Model, error) {
	var dbModel workerModel

	found, err := gorpmapping.Get(ctx, db, q, &dbModel)
	if err != nil {
		return nil, sdk.WrapError(err, "cannot get worker model")
	}
	if !found {
		return nil, sdk.WithStack(sdk.ErrNotFound)
	}

	isValid, err := gorpmapping.CheckSignature(dbModel, dbModel.Signature)
	if err != nil {
		return nil, err
	}
	if !isValid {
		log.Error(ctx, "workermodel.get> worker model %d data corrupted", dbModel.ID)
		return nil, sdk.WithStack(sdk.ErrNotFound)
	}

	model := dbModel.Model
	for i := range opts {
		if err := opts[i](ctx, db, &model); err != nil {
			return nil, err
		}
	}
	return &model, nil
}

func get(ctx context.Context, db gorp.SqlExecutor, q gorpmapping.Query, opts ...LoadOptionFunc) (*sdk.Model, error) {
	m, err := getCommon(ctx, db, q, opts...)
	if err != nil {
		return nil, err
	}

	// Temporary hide password
	if m.ModelDocker.Password != "" {
		m.ModelDocker.Password = sdk.PasswordPlaceholder
	}

	return m, nil
}

func getWithClearPassword(ctx context.Context, db gorp.SqlExecutor, q gorpmapping.Query, opts ...LoadOptionFunc) (*sdk.Model, error) {
	m, err := getCommon(ctx, db, q, opts...)
	if err != nil {
		return nil, err
	}

	// Temporary decrypt password
	if m.ModelDocker.Private && m.ModelDocker.Password != "" {
		m.ModelDocker.Password, err = secret.DecryptValue(m.ModelDocker.Password)
		if err != nil {
			return nil, sdk.WrapError(err, "cannot decrypt value for model with id %d", m.ID)
		}
	}

	return m, nil
}

// LoadAll retrieves worker models from database.
func LoadAll(ctx context.Context, db gorp.SqlExecutor, filter *LoadFilter, opts ...LoadOptionFunc) ([]sdk.Model, error) {
	var query gorpmapping.Query

	if filter == nil {
		query = gorpmapping.NewQuery("SELECT * FROM worker_model ORDER BY name")
	} else {
		query = gorpmapping.NewQuery(`
      SELECT distinct worker_model.*
      FROM worker_model
      LEFT JOIN worker_capability ON worker_model.id = worker_capability.worker_model_id
      WHERE ` + filter.SQL() + `
      ORDER BY worker_model.name
    `).Args(filter.Args())
	}

	return getAll(ctx, db, query, opts...)
}

// LoadAllByGroupIDs returns worker models list for given group ids.
func LoadAllByGroupIDs(ctx context.Context, db gorp.SqlExecutor, groupIDs []int64, filter *LoadFilter, opts ...LoadOptionFunc) ([]sdk.Model, error) {
	var query gorpmapping.Query

	if filter == nil {
		query = gorpmapping.NewQuery(`
      SELECT *
      FROM worker_model
      WHERE group_id = ANY($1)
      ORDER BY name
    `).Args(pq.Int64Array(groupIDs))
	} else {
		query = gorpmapping.NewQuery(`
      SELECT distinct worker_model.*
      FROM worker_model
      LEFT JOIN worker_capability ON worker_model.id = worker_capability.worker_model_id
      WHERE worker_model.group_id = ANY(:groupIDs)
      AND ` + filter.SQL() + `
      ORDER BY worker_model.name
    `).Args(filter.Args().Merge(gorpmapping.ArgsMap{
			"groupIDs": pq.Int64Array(groupIDs),
		}))
	}

	return getAll(ctx, db, query, opts...)
}

// LoadAllByNameAndGroupIDs retrieves all worker model with given name for group ids in database.
func LoadAllByNameAndGroupIDs(ctx context.Context, db gorp.SqlExecutor, name string, groupIDs []int64, opts ...LoadOptionFunc) ([]sdk.Model, error) {
	query := gorpmapping.NewQuery(`
    SELECT *
    FROM worker_model
    WHERE name = $1
    AND group_id = ANY($2)
    ORDER BY name
  `).Args(name, pq.Int64Array(groupIDs))
	return getAll(ctx, db, query, opts...)
}

// LoadAllActiveAndNotDeprecatedForGroupIDs retrieves models for given group ids.
func LoadAllActiveAndNotDeprecatedForGroupIDs(ctx context.Context, db gorp.SqlExecutor, groupIDs []int64, opts ...LoadOptionFunc) ([]sdk.Model, error) {
	query := gorpmapping.NewQuery(`
    SELECT *
    FROM worker_model
    WHERE group_id = ANY($1)
    AND is_deprecated = false
    AND disabled = false
    ORDER BY name
  `).Args(pq.Int64Array(groupIDs))
	return getAll(ctx, db, query, opts...)
}

// LoadByID retrieves a specific worker model in database.
func LoadByID(ctx context.Context, db gorp.SqlExecutor, id int64, opts ...LoadOptionFunc) (*sdk.Model, error) {
	query := gorpmapping.NewQuery(`
    SELECT *
    FROM worker_model
    WHERE id = $1
  `).Args(id)
	return get(ctx, db, query, opts...)
}

// LoadByNameAndGroupID retrieves a specific worker model in database by name and group id.
func LoadByNameAndGroupID(ctx context.Context, db gorp.SqlExecutor, name string, groupID int64, opts ...LoadOptionFunc) (*sdk.Model, error) {
	query := gorpmapping.NewQuery(`
    SELECT *
    FROM worker_model
    WHERE name = $1 AND group_id = $2
  `).Args(name, groupID)
	return get(ctx, db, query, opts...)
}

// LoadByIDWithClearPassword retrieves a specific worker model in database.
func LoadByIDWithClearPassword(ctx context.Context, db gorp.SqlExecutor, id int64, opts ...LoadOptionFunc) (*sdk.Model, error) {
	query := gorpmapping.NewQuery(`
    SELECT *
    FROM worker_model
    WHERE id = $1
  `).Args(id)
	return getWithClearPassword(ctx, db, query, opts...)
}

// LoadByNameAndGroupIDWithClearPassword retrieves a specific worker model in database by name and group id.
func LoadByNameAndGroupIDWithClearPassword(ctx context.Context, db gorp.SqlExecutor, name string, groupID int64, opts ...LoadOptionFunc) (*sdk.Model, error) {
	query := gorpmapping.NewQuery(`
    SELECT *
    FROM worker_model
    WHERE name = $1 AND group_id = $2
  `).Args(name, groupID)
	return getWithClearPassword(ctx, db, query, opts...)
}

// LoadAllUsableWithClearPasswordByGroupIDs returns usable worker models for given group ids.
func LoadAllUsableWithClearPasswordByGroupIDs(ctx context.Context, db gorp.SqlExecutor, groupIDs []int64, opts ...LoadOptionFunc) ([]sdk.Model, error) {
	// note about restricted field on worker model:
	// if restricted = true, worker model can be launched by a group hatchery only
	// so, a 'shared.infra' hatchery need all its worker models and all others with restricted = false

	query := gorpmapping.NewQuery(`
    SELECT *
    FROM worker_model
    WHERE (
      group_id = ANY($1)
      OR (
        $2 = ANY($1)
        AND restricted = false
      )
    ) AND disabled = false
    ORDER BY name
  `).Args(pq.Int64Array(groupIDs), group.SharedInfraGroup.ID)

	return getAllWithClearPassword(ctx, db, query, opts...)
}

// Insert a new worker model in database.
func Insert(ctx context.Context, db gorp.SqlExecutor, model *sdk.Model) error {
	dbmodel := workerModel{Model: *model}

	dbmodel.UserLastModified = time.Now()
	dbmodel.NeedRegistration = true

	if dbmodel.Type == sdk.Docker {
		dbmodel.ModelDocker.Envs = MergeModelEnvsWithDefaultEnvs(dbmodel.ModelDocker.Envs)
		if dbmodel.ModelDocker.Password == sdk.PasswordPlaceholder {
			return sdk.WithStack(sdk.ErrInvalidPassword)
		}
		if dbmodel.ModelDocker.Private {
			if dbmodel.ModelDocker.Password != "" {
				var err error
				dbmodel.ModelDocker.Password, err = secret.EncryptValue(dbmodel.ModelDocker.Password)
				if err != nil {
					return sdk.WrapError(err, "cannot encrypt docker password")
				}
			}
		} else {
			dbmodel.ModelDocker.Username = ""
			dbmodel.ModelDocker.Password = ""
			dbmodel.ModelDocker.Registry = ""
		}
	}

	if err := gorpmapping.InsertAndSign(ctx, db, &dbmodel); err != nil {
		return sdk.WithStack(err)
	}

	for _, r := range dbmodel.RegisteredCapabilities {
		if err := InsertCapabilityForModelID(db, dbmodel.ID, &r); err != nil {
			return err
		}
	}

	*model = dbmodel.Model

	// Temporary hide password
	if model.ModelDocker.Password != "" {
		model.ModelDocker.Password = sdk.PasswordPlaceholder
	}

	return nil
}

// UpdateDB a worker model
// if the worker model have SpawnErr -> clear them.
func UpdateDB(ctx context.Context, db gorp.SqlExecutor, model *sdk.Model) error {
	dbmodel := workerModel{Model: *model}

	if err := DeleteCapabilitiesByModelID(db, dbmodel.ID); err != nil {
		return err
	}

	dbmodel.UserLastModified = time.Now()
	dbmodel.NeedRegistration = true
	dbmodel.NbSpawnErr = 0
	dbmodel.LastSpawnErr = nil
	dbmodel.LastSpawnErrLogs = nil

	if dbmodel.ModelDocker.Password == sdk.PasswordPlaceholder {
		return sdk.WithStack(sdk.ErrInvalidPassword)
	}
	if dbmodel.ModelDocker.Private {
		if dbmodel.ModelDocker.Password != "" {
			var err error
			dbmodel.ModelDocker.Password, err = secret.EncryptValue(dbmodel.ModelDocker.Password)
			if err != nil {
				return sdk.WrapError(err, "cannot encrypt docker password")
			}
		}
	} else {
		dbmodel.ModelDocker.Username = ""
		dbmodel.ModelDocker.Password = ""
		dbmodel.ModelDocker.Registry = ""
	}

	if err := gorpmapping.UpdateAndSign(ctx, db, &dbmodel); err != nil {
		return sdk.WithStack(err)
	}

	for _, r := range dbmodel.RegisteredCapabilities {
		if err := InsertCapabilityForModelID(db, dbmodel.ID, &r); err != nil {
			return err
		}
	}

	*model = dbmodel.Model

	// Temporary hide password
	if model.ModelDocker.Password != "" {
		model.ModelDocker.Password = sdk.PasswordPlaceholder
	}

	return nil
}

// DeleteByID a worker model from database and all its capabilities.
func DeleteByID(db gorp.SqlExecutor, id int64) error {
	_, err := db.Exec("DELETE FROM worker_model WHERE id = $1", id)
	return sdk.WrapError(err, "unable to remove worker model with id %d", id)
}
