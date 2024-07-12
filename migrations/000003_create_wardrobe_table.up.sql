CREATE TABLE Wardrobe (
                          id SERIAL PRIMARY KEY,
                          Title VARCHAR(255) NOT NULL,
                          Quantity INT NOT NULL,
                          Price VARCHAR(255) NOT NULL,
                          Old_price VARCHAR(225) NOT NULL,
                          Description TEXT,
                          Height VARCHAR(255),
                          Width VARCHAR(255),
                          Depth VARCHAR(255),
                          Filename VARCHAR(255),
                          Link VARCHAR(255)
);