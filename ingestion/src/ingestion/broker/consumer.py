import json
from collections.abc import Callable
from typing import Any

from confluent_kafka import Consumer


class KafkaConsumer:
    def __init__(self, bootstrap_servers: str, group_id: str, topics: list[str]):
        self.consumer = Consumer({
            "bootstrap.servers": bootstrap_servers,
            "group.id": group_id,
            "auto.offset.reset": "earliest",
            "enable.auto.commit": False,
        })

        self.consumer.subscribe(topics)

    def consume(self, handler: Callable[[dict[str, Any]], None]) -> None:
        try:
            while True:
                message = self.consumer.poll(1.0)

                if message is None: continue
                if message.error(): raise RuntimeError(message.error())

                raw_value = message.value()

                if raw_value is None:
                    self.consumer.commit(message=message)
                    continue

                value = json.loads(raw_value.decode("utf-8"))

                handler(value)

                self.consumer.commit(message=message)

        finally:
            self.consumer.close()
