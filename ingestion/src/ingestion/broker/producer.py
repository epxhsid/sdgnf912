import json
from typing import Any

from confluent_kafka import Producer


class KafkaProducer:
    def __init__(self, bootstrap_servers: str):
        self.producer = Producer({
            "bootstrap.servers": bootstrap_servers
        })

    def produce(self, topic: str, value: dict[str, Any], key: str | None = None) -> None:
        self.producer.produce(
            topic=topic,
            key=key,
            value=json.dumps(value).encode("utf-8"),
        )

    def flush(self) -> None:
        self.producer.flush()
