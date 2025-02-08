DO $$ 
BEGIN

    CREATE TABLE IF NOT EXISTS results (
	    id SERIAL PRIMARY KEY,
	    host VARCHAR(255) NOT NULL UNIQUE,
	    min_time FLOAT NOT NULL,
	    max_time FLOAT NOT NULL,
	    last_up  VARCHAR(255) DEFAULT 'never',
	    ping_time VARCHAR(255)
    );

    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        username VARCHAR(255) NOT NULL,
        password VARCHAR(255) NOT NULL
    );


    RAISE NOTICE 'Table "results" created successfully';
END $$;