from datetime import UTC, datetime, timedelta
from uuid import UUID

from sqlalchemy import select, update
from sqlalchemy.ext.asyncio import AsyncSession

from ingestion.database.models import OutboxEventModel
from ingestion.models.outbox import OutboxEvent


class OutboxPersistence:
    def __init__(self, session: AsyncSession):
        self.session = session

    async def create(self, event: OutboxEvent) -> None:
        payload = {
            **event.payload,
            "event_id": str(event.id),
        }

        model = OutboxEventModel(
            id=event.id,
            event_type=event.event_type,
            aggregate_id=event.aggregate_id,
            payload=payload,
        )

        self.session.add(model)

    async def get_pending(self, limit: int = 100) -> list[OutboxEventModel]:
        result = await self.session.execute(
            select(OutboxEventModel)
            .where(
                OutboxEventModel.published_at.is_(None),
                OutboxEventModel.next_retry_at <= datetime.now(UTC),
            )
            .order_by(OutboxEventModel.created_at)
            .limit(limit)
        )

        return list(result.scalars())

    async def mark_published(self, event_id: UUID) -> None:
        await self.session.execute(
            update(OutboxEventModel)
            .where(OutboxEventModel.id == event_id)
            .values(
                published_at=datetime.now(UTC),
            )
        )

    async def mark_failed(self, event_id: UUID, error: str) -> None:
        now = datetime.now(UTC)
        await self.session.execute(
            update(OutboxEventModel)
            .where(OutboxEventModel.id == event_id)
            .values(
                attempts=OutboxEventModel.attempts + 1,
                last_error=error,
                next_retry_at=now + timedelta(seconds=5),
            )
        )
