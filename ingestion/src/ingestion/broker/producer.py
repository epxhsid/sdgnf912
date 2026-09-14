import json
from typing import Any

from confluent_kafka.aio import AIOProducer


class KafkaProducer:
    def __init__(self, bootstrap_servers: str):
        self.producer = AIOProducer(
            {
                "bootstrap.servers": bootstrap_servers,
                "message.timeout.ms": 5000,
            }
        )

    async def produce(
        self,
        topic: str,
        value: dict[str, Any],
        key: str | None = None,
    ):
        return await self.producer.produce(
            topic=topic,
            key=key,
            value=json.dumps(value).encode("utf-8"),
        )

    async def flush(self) -> None:
        await self.producer.flush()

    async def close(self) -> None:
        await self.producer.close()
