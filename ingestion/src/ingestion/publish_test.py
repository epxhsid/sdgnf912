import asyncio
import logging

from ingestion.broker.producer import KafkaProducer
from ingestion.broker.publisher import OutboxPublisher
from ingestion.config import KAFKA_BOOTSTRAP_SERVERS
from ingestion.database.engine import session_factory

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s %(levelname)s %(name)s: %(message)s",
)

async def main() -> None:
    producer = KafkaProducer(
        bootstrap_servers=KAFKA_BOOTSTRAP_SERVERS,
    )

    publisher = OutboxPublisher(
        session_factory=session_factory,
        producer=producer,
    )

    await publisher.run()

if __name__ == "__main__":
    asyncio.run(main())
