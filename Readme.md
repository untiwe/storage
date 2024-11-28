### Учебная программа

Собрать dockerfile
`docker build --tag storage ./main`
Запустить docker pompose
`docker-compose up`
Готово. 

Я попотылтся мнимизировать настройки перед первым запуском. БД и Кафка собираются сами с нуля, но докер собрать всё же придется.

По адресу `http://localhost:8080/` будет находится страница с UI ее использование показано на видео
Для создания заканов нужен POST запрос `http://localhost:8080/addorder/`
JSON как в задаче
```{
	"order_uid": "wbschool-94d6-468b-9024-01e10a8da71e",
	"track_number": "TRACK970388",
	"entry": "wjeV",
	"delivery": {
		"name": "iYt3ZiDwcr",
		"phone": "+9591010157",
		"zip": "307893",
		"city": "OfutXz3U2l",
		"address": "aT7oYU0uP6sk8ivspQOz",
		"region": "kYYmIvkIUB",
		"email": "Rg2iM@byWh5.com"
	},
	"payment": {
		"transaction": "TRANS32781",
		"request_id": "MBYp9gqZB4",
		"currency": "TUo",
		"provider": "cthMgEGsUX",
		"amount": 6987,
		"payment_dt": 1732793571,
		"bank": "4nscPMBRCk",
		"delivery_cost": 558,
		"goods_total": 741,
		"custom_fee": 55
	},
	"items": [
		{
			"chrt_id": 172236,
			"track_number": "TRACK962598",
			"price": 775,
			"rid": "RID416796",
			"name": "Ax8sSy3iHi",
			"sale": 78,
			"size": "Yj",
			"total_price": 245,
			"nm_id": 86924,
			"brand": "kdlIVp3Sig",
			"status": 6
		}
	],
	"locale": "Pz",
	"internal_signature": "HcuhaWktGs",
	"customer_id": "TKQFDduEmX",
	"delivery_service": "n4EDLCy7pX",
	"shardkey": "b",
	"sm_id": 95,
	"date_created": "2024-11-28T11:32:51.000743138Z",
	"oof_shard": "f"
}
```

Урлы и порты специально не выносил в конфиги т.к. в этом проекте они меняться никогда не будут, как и мнодество других параметров. Вынес то, что было в задаче или может хоть потенциально может прогодится. При желании можно прокинуть файл с наружи, я закомеентировл этот volume в docker compose. 

Ещё раз дублирую ссылку на видео (в бенде) https://band.wb.ru/files/rq56mhjo9bbizpsjxe5gptmere/public?h=M2gvGRVdTePSMxMZfftzkjIdHpqLYATTDV50I6FPiRk