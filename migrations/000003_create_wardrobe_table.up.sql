CREATE TABLE Wardrobe (
                          id SERIAL PRIMARY KEY,
                          Title VARCHAR(255) NOT NULL,
                          Quantity INT NOT NULL,
                          Price DECIMAL(10, 2) NOT NULL,
                          Old_price DECIMAL(10, 2) NOT NULL,
                          Description TEXT,
                          Height DECIMAL(10, 2),
                          Width DECIMAL(10, 2),
                          Depth DECIMAL(10, 2),
                          Filename VARCHAR(255),
                          Link VARCHAR(255)
);