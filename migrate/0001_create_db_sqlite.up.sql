
CREATE TABLE IF NOT EXISTS tb_product (
    id         INTEGER         PRIMARY KEY AUTOINCREMENT UNIQUE NOT NULL,
    name        VARCHAR (120)   NOT NULL,
    code       VARCHAR (20)    NOT NULL,
    price       DECIMAL (10, 2) NOT NULL,
    created_at DATETIME        NOT NULL DEFAULT (datetime('now', 'localtime')),
    updated_at DATETIME        NOT NULL DEFAULT (datetime('now', 'localtime'))
);

CREATE TRIGGER IF NOT EXISTS tb_product_update_trig
         AFTER UPDATE OF name,
                         code,
                         price
            ON tb_product
      FOR EACH ROW
BEGIN
    UPDATE tb_product
       SET updated_at = datetime('now', 'localtime') 
     WHERE id = NEW.id;
END;
