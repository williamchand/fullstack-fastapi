"""migrate account and apps basic

Revision ID: 3623c47a4217
Revises: 1a31ce608336
Create Date: 2025-11-03 16:08:46.162080

"""
from alembic import op
import sqlalchemy as sa
import sqlmodel.sql.sqltypes


# revision identifiers, used by Alembic.
revision = '3623c47a4217'
down_revision = '1a31ce608336'
branch_labels = None
depends_on = None


def upgrade():
    pass


def downgrade():
    pass
