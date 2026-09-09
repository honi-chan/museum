package exhibition

import "context"

// CreateAndSave は、
// Exhibition を作成してRepositoryへ保存する。
//
// この処理は、
//
// 1. Exhibitionを作る
// 2. Exhibitionを保存する
//
// という一連のユースケースを担当する。
//
// Exhibitionの生成ルール自体は Create() に任せる。
func CreateAndSave(
	ctx context.Context,
	input CreateInput,
	idGenerator IDGenerator,
	repository Repository,
) (Exhibition, error) {
	// Exhibitionそのものを生成する。
	//
	// タイトル必須などのルールは、
	// Create()側に集約している。
	created, err := Create(
		input,
		idGenerator,
	)

	if err != nil {
		return Exhibition{}, err
	}

	// 保存方法については知らない。
	//
	// PostgreSQLなのかメモリなのかではなく、
	// Repository.Save()だけを利用する。
	if err := repository.Save(
		ctx,
		created,
	); err != nil {
		return Exhibition{}, err
	}

	return created, nil
}
