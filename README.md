# runtemp

## Para executar os serviços localmente via Docker, siga os passos abaixo:
Clone este repositório.
$ git clone https://github.com/dukivt/runtemp

Na raiz do projeto, execute o seguinte comando: **#docker-compose up -d**

```

## Testando Cenários
execute o arquivo requisicoes.http, que se encontra na raiz do projeto
ou
Acessar: http://localhost:8082/?cep={cep}
Ex: http://http://localhost:8082/?cep=24230136

```

## Acesso via hospedagem Google Cloud Run
Acessar o endereço abaixo informando o cep que deseja consultar:
https://temperaturacep-4tynn5krra-ue.a.run.app/?cep={cep}

Exemplo: https://temperaturacep-4tynn5krra-ue.a.run.app/?cep=24230136
