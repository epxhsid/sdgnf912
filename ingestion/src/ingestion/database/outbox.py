from ingestion.database.models import Base


class OutboxEventModel(Base):
    __tablename__ = "outbox_events"
