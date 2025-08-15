-- cria tabela para os parceiros visando autenticação
CREATE TABLE partners (
    id INT AUTO_INCREMENT PRIMARY KEY, 
    name VARCHAR(50) NOT NULL, 
    api_key VARCHAR(100) NOT NULL UNIQUE -- garantir chave diferetes
); 

-- cria tabela para conversões
CREATE TABLE conversions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    transaction_id VARCHAR(50) NOT NULL UNIQUE,  -- sem duplicatas
    partner_id INT NOT NULL, 
    amount DECIMAL(6, 2) NOT NULL, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    FOREIGN KEY (partner_id) REFERENCES partners(id)

);

INSERT INTO partners (name, api_key) VALUES ('ParceiroX', 'davi-chave-1234');