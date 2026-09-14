import asyncio
import logging

from ingestion.broker.producer import KafkaProducer
from ingestion.database.outbox import OutboxPersistence

logger = logging.getLogger(__name__)


class OutboxPublisher:
    def __init__(self, session_factory, producer: KafkaProducer, poll_interval: float = 1.0, batch_size: int = 100):
        self.session_factory = session_factory
        self.producer = producer
        self.poll_interval = poll_interval
        self.batch_size = batch_size
        self.running = False

    async def run(self) -> None:
        self.running = True

        logger.info("Outbox publisher started")

        try:
            while self.running:
                published = await self.publish_batch()

                if published == 0:
                    await asyncio.sleep(self.poll_interval)

        finally:
            logger.info("Flushing Kafka producer")

            remaining = self.producer.flush()

            if remaining:
                logger.warning(
                    "%d messages were not delivered",
                    remaining,
                )

            logger.info("Outbox publisher stopped")

    def stop(self) -> None:
        self.running = False

    async def publish_batch(self) -> int:
        async with self.session_factory() as session:
            outbox = OutboxPersistence(session)

            events = await outbox.get_pending(
                limit=self.batch_size,
            )

        if not events:
            return 0

        published = 0

        for event in events:
            success = self._publish_event(event)

            if success:
                async with self.session_factory() as session:
                    outbox = OutboxPersistence(session)

                    await outbox.mark_published(event.id)

                    await session.commit()

                published += 1

        return published

    def _publish_event(self, event) -> bool:
        delivery_error = None

        def on_delivery(error, message):
            nonlocal delivery_error

            if error is not None:
                delivery_error = error
                return

            logger.info(
                "Published event %s to %s [%s] @ %s",
                event.id,
                message.topic(),
                message.partition(),
                message.offset(),
            )

        try:
            self.producer.produce(
                topic=event.event_type,
                key=str(event.aggregate_id),
                value=event.payload,
                on_delivery=on_delivery,
            )

            remaining = self.producer.flush()

            if delivery_error is not None:
                raise delivery_error

            if remaining:
                raise RuntimeError(
                    f"{remaining} Kafka messages remain undelivered"
                )

            return True

        except Exception:
            logger.exception(
                "Failed to publish event %s",
                event.id,
            )

            return False
