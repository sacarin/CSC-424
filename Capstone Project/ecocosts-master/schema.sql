CREATE TABLE client (
	id SERIAL,
	name TEXT UNIQUE NOT NULL,
	pass TEXT NOT NULL,
	PRIMARY KEY (id),
	CHECK (LENGTH(name) < 20)
);

CREATE TABLE category (
	id SERIAL,
	description TEXT NOT NULL,
	PRIMARY KEY (id)
);

CREATE TABLE budget (
	client_id INT NOT NULL,
	cat_id INT NOT NULL,
	amount MONEY NOT NULL,
	FOREIGN KEY(client_id) REFERENCES client(id) ON DELETE CASCADE,
	FOREIGN KEY(cat_id) REFERENCES category(id)
);

CREATE TABLE transaction (
	id SERIAL,
	client_id INT NOT NULL,
	cat_id INT,
	amount MONEY NOT NULL,
	balance MONEY NOT NULL DEFAULT 0.00,
	description TEXT,
	time TIMESTAMP NOT NULL DEFAULT 'now'::timestamp,
	PRIMARY KEY (id),
	FOREIGN KEY (client_id) REFERENCES client(id) ON DELETE CASCADE,
	FOREIGN KEY (cat_id) REFERENCES category(id)
);

CREATE TABLE stock (
	client_id INT NOT NULL,
	symbol TEXT NOT NULL,
	quantity INT NOT NULL,
	FOREIGN KEY (client_id) REFERENCES client(id) ON DELETE CASCADE,
	CHECK (quantity > 0)
);	

CREATE FUNCTION calc_balance() RETURNS trigger AS $$
BEGIN
	UPDATE transaction
	SET balance = coalesce(NEW.amount + (
		SELECT DISTINCT balance
		FROM transaction
		WHERE client_id = NEW.client_id AND time = (
			SELECT MAX(time)
			FROM transaction
			WHERE client_id = NEW.client_id AND time < NEW.time
		)
	), NEW.amount)
	WHERE time = NEW.time AND client_id = NEW.client_id;
RETURN NULL;
END;
$$ LANGUAGE plpgsql VOLATILE COST 100;

CREATE TRIGGER update_calc_balance
AFTER INSERT OR UPDATE OF amount
ON transaction
FOR EACH ROW
EXECUTE PROCEDURE calc_balance();
