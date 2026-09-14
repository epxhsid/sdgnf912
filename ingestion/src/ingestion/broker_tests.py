from ingestion.broker.consumer import KafkaConsumer
from ingestion.broker.producer import KafkaProducer

BOOTSTRAP_SERVERS = "localhost:19092"


def handle_message(message: dict) -> None:
    print("Received:", message)


def main() -> None:
    producer = KafkaProducer(BOOTSTRAP_SERVERS)

    producer.produce(
        topic="document.uploaded",
        key="test-document",
        value={
            "event_id": "test-event-1",
            "event_type": "document.uploaded",
            "document_id": "test-document-1",
            "source_uri": "seaweedfs://documents/test.pdf",
            "filename": "test.pdf",
            "content_type": "application/pdf",
        },
    )

    producer.flush()

    consumer = KafkaConsumer(
        bootstrap_servers=BOOTSTRAP_SERVERS,
        group_id="argus-ingestion-test",
        topics=["document.uploaded"],
    )

    consumer.consume(handle_message)


if __name__ == "__main__":
    main()
