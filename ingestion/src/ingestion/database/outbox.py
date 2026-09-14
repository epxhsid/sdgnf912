from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession

from ingestion.database.models import OutboxEventModel
from ingestion.models.outbox import OutboxEvent


class OutboxPersistence:
    def __init__(self, session: AsyncSession):
        self.session = session

    async def create(self, event: OutboxEvent) -> None:
        model = OutboxEventModel(
            id=event.id,
            event_type=event.event_type,
            aggregate_id=event.aggregate_id,
            payload=event.payload,
        )

        self.session.add(model)

    async def get_pending(self, limit: int = 100) -> list[OutboxEvent]:
        result = await self.session.execute(
            select(OutboxEventModel)
            .where(OutboxEventModel.published_at.is_(None))
            .order_by(OutboxEventModel.created_at)
            .limit(limit)
        )

        return [
            OutboxEvent(
                id=model.id,
                event_type=model.event_type,
                aggregate_id=model.aggregate_id,
                payload=model.payload,
            )
            for model in result.scalars()
        ]
