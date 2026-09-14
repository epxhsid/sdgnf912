import json
from collections.abc import Callable
from typing import Any

from confluent_kafka import Producer

DeliveryCallback = Callable[[Any, Any], None]

class KafkaProducer:
    def __init__(self, bootstrap_servers: str):
        self.producer = Producer(
            {
                "bootstrap.servers": bootstrap_servers,
                "message.timeout.ms": 5000,
            }
        )

    def produce(self, topic: str, value: dict[str, Any], key: str | None = None, on_delivery: DeliveryCallback | None = None) -> None:
        self.producer.produce(
            topic=topic,
            key=key,
            value=json.dumps(value).encode("utf-8"),
            callback=on_delivery,
        )

    def poll(self, timeout: float = 0.0) -> None:
        self.producer.poll(timeout)

    def flush(self, timeout: float = 10.0) -> int:
        return self.producer.flush(timeout)
