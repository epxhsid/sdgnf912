from uuid import UUID, uuid4

from pydantic import BaseModel, Field


class OutboxEvent(BaseModel):
    id: UUID = Field(default_factory=uuid4)
    event_type: str
    aggregate_id: UUID
    payload: dict
