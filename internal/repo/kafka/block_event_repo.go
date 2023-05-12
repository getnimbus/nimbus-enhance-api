package kafka

import (
	"context"

	"go-nimeth/internal/conf"
	"go-nimeth/internal/entity_dto"
	"go-nimeth/internal/infra"
	"go-nimeth/internal/repo"
)

func NewBlockEventRepo(kafkaProducer infra.KafkaSyncProducer) repo.BlockEventRepo {
	return &blockEventRepo{
		topic:         conf.Config.NewBlockEventsTopic,
		kafkaProducer: kafkaProducer,
	}
}

type blockEventRepo struct {
	topic         string
	kafkaProducer infra.KafkaSyncProducer
}

func (repo *blockEventRepo) Send(ctx context.Context, msg *entity_dto.BlockEvent) error {
	return repo.kafkaProducer.SendJson(ctx, repo.topic, msg)
}
