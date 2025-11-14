"""Initialize wedding and payment models

Revision ID: 9c0a54914c78
Revises: e2412789c190
Create Date: 2024-06-17 14:42:44.639457

"""
from alembic import op
import sqlalchemy as sa
import sqlalchemy.dialects.postgresql as pg


# revision identifiers, used by Alembic.
revision = '9c0a54914c78'
down_revision = 'e2412789c190'
branch_labels = None
depends_on = None


def upgrade():
    # --- payment_method table ---
    op.create_table(
        'payment_method',
        sa.Column('id', pg.UUID(as_uuid=True), primary_key=True, server_default=sa.text('gen_random_uuid()')),
        sa.Column('name', sa.String(length=100), nullable=False, unique=True),
        sa.Column('provider', sa.String(length=50), nullable=False),  # e.g. stripe, qris, bank_transfer
        sa.Column('config', pg.JSONB(), nullable=False, server_default=sa.text("'{}'::jsonb")),
        sa.Column('is_active', sa.Boolean(), nullable=False, server_default=sa.text('true')),
        sa.Column('created_at', sa.DateTime(timezone=True), server_default=sa.text('NOW()')),
    )

    # --- payment table ---
    op.create_table(
        'payment',
        sa.Column('id', pg.UUID(as_uuid=True), primary_key=True, server_default=sa.text('gen_random_uuid()')),
        sa.Column('user_id', pg.UUID(as_uuid=True), sa.ForeignKey('user.id', ondelete='CASCADE'), nullable=False),
        sa.Column('payment_method_id', pg.UUID(as_uuid=True), sa.ForeignKey('payment_method.id', ondelete='SET NULL'), nullable=True),
        sa.Column('amount', sa.Numeric(10, 2), nullable=False),
        sa.Column('currency', sa.String(length=10), nullable=False, server_default='USD'),
        sa.Column('status', sa.Enum('pending', 'paid', 'failed', name='payment_status'), nullable=False, server_default='pending'),
        sa.Column('transaction_id', sa.String(length=255), nullable=False, unique=True),
        sa.Column('metadata', pg.JSONB(), nullable=False, server_default=sa.text("'{}'::jsonb")),
        sa.Column('created_at', sa.DateTime(timezone=True), server_default=sa.text('NOW()')),
    )

    # --- template table ---
    op.create_table(
        'template',
        sa.Column('id', pg.UUID(as_uuid=True), primary_key=True, server_default=sa.text('gen_random_uuid()')),
        sa.Column('name', sa.String(length=100), nullable=False, unique=True),
        sa.Column('theme_config', pg.JSONB(), nullable=False, server_default=sa.text("'{}'::jsonb")),
        sa.Column('config_schema', pg.JSONB(), nullable=False, server_default=sa.text("'{}'::jsonb")),
        sa.Column('preview_url', sa.String(length=512), nullable=True),
        sa.Column('price', sa.Numeric(10, 2), nullable=False, server_default='0.00'),
        sa.Column('created_at', sa.DateTime(timezone=True), server_default=sa.text('NOW()')),
    )

    # --- wedding table ---
    op.create_table(
        'wedding',
        sa.Column('id', pg.UUID(as_uuid=True), primary_key=True, server_default=sa.text('gen_random_uuid()')),
        sa.Column('user_id', pg.UUID(as_uuid=True), sa.ForeignKey('user.id', ondelete='CASCADE'), nullable=False),
        sa.Column('template_id', pg.UUID(as_uuid=True), sa.ForeignKey('template.id', ondelete='SET NULL'), nullable=True),
        sa.Column('payment_id', pg.UUID(as_uuid=True), sa.ForeignKey('payment.id', ondelete='SET NULL'), nullable=True),
        sa.Column('status', sa.Enum('draft', 'published', name='wedding_status'), nullable=False, server_default='draft'),
        sa.Column('custom_domain', sa.String(length=255), nullable=True, unique=True),
        sa.Column('slug', sa.String(length=150), nullable=True),
        sa.Column('config_data', pg.JSONB(), nullable=False, server_default=sa.text("'{}'::jsonb")),
        sa.Column('created_at', sa.DateTime(timezone=True), server_default=sa.text('NOW()')),
        sa.Column('deleted_at', sa.DateTime(timezone=True), nullable=True), # soft delete
    )
    op.create_index(
        'ix_wedding_slug_active',
        'wedding',
        ['slug'],
        unique=True,
        postgresql_where=sa.text('deleted_at IS NULL')
    )

    # --- guest table ---
    op.create_table(
        'guest',
        sa.Column('id', pg.UUID(as_uuid=True), primary_key=True, server_default=sa.text('gen_random_uuid()')),
        sa.Column('wedding_id', pg.UUID(as_uuid=True), sa.ForeignKey('wedding.id', ondelete='CASCADE'), nullable=False),
        sa.Column('name', sa.String(length=255), nullable=False),
        sa.Column('contact', sa.String(length=255), nullable=False),
        sa.Column('rsvp_status', sa.Enum('yes', 'no', 'maybe', name='rsvp_status'), nullable=False, server_default='maybe'),
        sa.Column('message', sa.Text(), nullable=True),
        sa.Column('created_at', sa.DateTime(timezone=True), server_default=sa.text('NOW()')),
        sa.Column('deleted_at', sa.DateTime(timezone=True), nullable=True),   # soft delete
    )
    # Replace unique constraint with partial unique index
    op.create_index(
        "uix_guest_wedding_contact_active",
        "guest",
        ["wedding_id", "contact"],
        unique=True,
        postgresql_where=sa.text("deleted_at IS NULL")
    )

    # --- optional seed ---
    op.execute("""
        INSERT INTO payment_method (name, provider, config)
        VALUES
            ('Stripe', 'stripe', '{}'),
            ('QRIS', 'qris', '{}'),
            ('Manual Bank Transfer', 'bank_transfer', '{}')
        ON CONFLICT (name) DO NOTHING;
    """)


def downgrade():
    op.drop_table('guest')
    op.drop_table('wedding')
    op.drop_table('template')
    op.drop_table('payment')
    op.drop_table('payment_method')

    op.execute('DROP TYPE IF EXISTS wedding_status CASCADE;')
    op.execute('DROP TYPE IF EXISTS rsvp_status CASCADE;')
    op.execute('DROP TYPE IF EXISTS payment_status CASCADE;')