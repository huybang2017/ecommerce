-- Initialize service databases if they do not already exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'identity_service') THEN
        CREATE DATABASE identity_service;
    END IF;

    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'product_service') THEN
        CREATE DATABASE product_service;
    END IF;

    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'order_service') THEN
        CREATE DATABASE order_service;
    END IF;

    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'inventory_service') THEN
        CREATE DATABASE inventory_service;
    END IF;

    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'promotion_service') THEN
        CREATE DATABASE promotion_service;
    END IF;

    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'payment_service') THEN
        CREATE DATABASE payment_service;
    END IF;

    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'notification_service') THEN
        CREATE DATABASE notification_service;
    END IF;
END
$$;
