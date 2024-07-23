# Go Expert

Desafio **Sistema de Temperatura por CEP** do curso **Pós Go Expert**.

**Objetivo:** Desenvolver um sistema em Go que receba um CEP, identifica a cidade e retorna o clima atual (temperatura em graus celsius, fahrenheit e kelvin). Esse sistema deverá ser publicado no Google Cloud Run.

### Execução da **aplicação**
Para executar a aplicação execute o comando:
```
git clone https://github.com/IgorLopes88/goexpert-cloudrun.git
cd goexpert-cloudrun
go mod tidy
go run main.go
```

Em seguida abra o navegador e acesse o endereço `http://localhost:8080/temperature/00000-000`, onde `00000-000` deverá ser subistituido pelo **CEP DESEJADO**.

Exemplo:

```
http://localhost:8080/temperature/11600300
```

### Execução da **aplicação** via Docker
Para executar a aplicação execute o comando:
```
git clone https://github.com/IgorLopes88/goexpert-cloudrun.git
cd goexpert-cloudrun
go mod tidy
docker-compose up
```

Em seguida abra o navegador e acesse o endereço `http://localhost:8080/temperature/00000-000`, onde `00000-000` deverá ser subistituido pelo **CEP DESEJADO**. 

Exemplo:

```
http://localhost:8080/temperature/11666000
```
### Execução da **aplicação** via Cloud Run
Para acessar a aplicação abra o navegador e acesse o endereço `https://cloud-run-goexpert-h5qxusn7mq-uc.a.run.app/temperature/00000-000`, onde `00000-000` deverá ser subistituido pelo **CEP DESEJADO**. 

Exemplo:

```
https://cloud-run-goexpert-h5qxusn7mq-uc.a.run.app/temperature/11633000
```


Pronto!


### Correções de Bugs
