package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"math/big"
	"sync"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"golang.org/x/crypto/sha3"
)

const targetByte = 0x01 // Задаем значение для проверки первого байта

// Список операторских ключей
var operatorList = []string{
	"NeonPQFrw5stVvs1rFLDxALWUBDCnSPsWBP83RfNUKK",
	"GYt9w8MaXztDLhhsxmQr7Ar9FJ6MmaFwav7qBrxZKwhd",
	"65g8N1ZTkWNo953PuWcYvnNsqC7EmbDvtRtgp4bmVT4S",
	"EcRUp3ah4CfAkew6UALfjHdzxZNQkWo767YTKK4H5HhM",
	"AbCp2zzd3qcA14uBosLdaYWUEJMxMBb2sfkzhDWURszW",
	"5qT6Wh3FyY3jJdqcf39ZXXMRN3zc2wEqQAKuakkXBZdG",
	"GfvDvxwjngmDNHfy6vpSk7MsRXzcfFtjCPe9jFuRJLih",
	"DnM2coxi7qu3AEKhN6Qk9mGQEAL1mAzeRmQyyWvQhcNX",
	"GJi4s1ALwrcGrdgAEZ6o2H6w4Movs2KwtFmZfZmt3eDN",
	"5atJcg4d8TAqAjRynWXHCsPs7gbJgqWHYvVkERMnkRQs",
	"EWhrrtA5fv9YtcHPWmYMK8VBS1V1t5s4GTspxPPuXGPZ",
	"45wTPXC8VvqrW3rymWVmrWsF4LisNnzJTnniU9LnRHJH",
	"CbdwxL2txJCGnYbEAMohfdcm7fUZWL87aCLjsqSpyKzR",
	"273fHsaWQwDcXEuzJ9yS2wu4VvZsRxc5fuK2oSmNY513",
	"DC1M1eD1ZYdUYWVy1GKr372p1H7jJVTSAjzmgaPgFSHW",
	"6KZLqnACXZrZ5JuPCm55m7YWesTNnCbGGDr3RysxhhCS",
	"ACzzSM2caeYueQCn9x8HpMaqoEkYFiigxikaxqe5iNcE",
	"5S7onsbfUXRMgrDX7R2Lp7v24z21kBoRVKeyEnAjGcdL",
	"GFa3c88L9NdkSKnwNyBCUC6juHSXPifat5r1u7PxZQQD",
	"EybjaSYmBcAzvfEtPk3fL1xtrTLvBSYAZ3CuSzu2nXrp",
	"HhDDQWtnDufThqK2tf1YfYbZzSouhErtuzRx8y4thsei",
	"GdkdE3husHhFzdt37MzWh8HPBi2MXRJHCZHkSCNLfTuD",
	"U1opJYxJV4zuQKuFwoXW45dLRKpNo7EkJjguRqWvy9s",
	"BAkYYQNNziqW2HmxQDSCU57JmqW2cmmfdeqmduhk5gYm",
	"KKjaYCyjXAfwuugaMSpgE2widgB6GUajD5Pze8NQptm",
	"EpuLUGJdSA9Lhe6pNxfXEAfBo6eGEuTVJjRQY93ZHdqP",
	"5i2XZkeLZuQKuns8aKDAKMnp2UCtKU9H4tYT4npXax4g",
	"DHGV5vXk9HpE7BRbCPmruFdq3wXBHqYBq4b8AP5vD6pN",
	"BhiWQ67L79bqAvefzEw9en3hmodXdyCR5a5qgAPT7EQU",
	"9GprgQo6pJ4hhT9YGAUeFsViE7YRaK8dttswYeN1tN6G",
	"AvCNVjRx8m9mycvhn87mbVbbi3YrgVJjV2fAFu9rsgk1",
	"HhbDiF6v8aPywr2gUU9AdQ2PX9W73ARc9bVQjEc931sm",
	"DJoMJKcksFEtjvjE7xkLKrT5kcVcniQR9s4C6ZuD3UNk",
	"HAFHnF4X2xkNjUTrXEw1FWKRF4kXdmn9wN61VA6kkuz7",
	"DAf6MQokeK7KdUVXCC6REX9DNwssAXEgt4Ri3vczT6HV",
	"2fNHBqeArUXSz32DfXQC4AmiiTNuo7mLRVWVY8tPZKKF",
	"9WYXXkQ8CSJd4t2Q3HT4auQVJyE7B86jodpqrDpw2XDE",
	"5p6SzmyEPBgySFSLZoSF5CXnUU7GPpf67f6NjbfhdKBH",
	"4NxuZhiHzeL6pENbPpP4GuHadHHu1jSDGkz5MDBW1tMR",
	"2uyrreutuU7s9FjgV3ptsX9UhNiboi1GwxXunXiVwmcH",
	"7VDZcJqZv2sUW5XmVtyGqwpFpTpbaz88ppGjyBbtUAZ3",
	"ALLcid9tKeSHumLZgmYoaFnPTUgTVo9i63quf65v5j6Y",
	"EPCPeqVgi8TP91zhYrzK9pz9mwnTsqxdH6VBXBnXoPjr",
	"AtvULAyvCMXSPeBa3oAM4y9J18ncW4VKMDtSibYrBA6T",
	"DL1qqtqcZrijd8WrfvwZMBe8stUCSntjVxWHfVF4dfjC",
	"Gct64jtD72WUFjgjd1Lz65e6Rc2Edkt4QuEkb82xujZd",
	"9cpMJQ5i1tDWTh4Ya2shJAXEM7AygHbYfgQvYTB13MX",
	"9hgaCEx6YepnXuJaRzQHoGovdx3bi988cvnQR48tX2xw",
	"EzCyipR6EKS494dWuz1hcxkoQR1b6dBAU7QGZCksYvvb",
	"qMpnDCegWrzx8MGy63ikvvgwQxuubmHza3QaGiLTgGL",
	"6fRgRzU2caLjeFGcA9brkZudwRc8xczN85NsEBVLrfXK",
	"ECER7q7cNcW1HTVpoKpfcNQh2NQ1WwusKvkk7CA5bYjB",
	"9db14BTMSmHLkUQBce9SB8fstLV84AxwaTnAS3JTFkmn",
	"aiRv8rmf4wXeQVoBRnCEcJEYoPzU4ovQVECo173VZeH",
	"GZmJa2gdvZWFiqEvnCthB5VBb5W5KmsbJEqSztUPxUwf",
	"ACpyy8i6q2fxncDvAbCBy7KMVkdkAgg6z6GfB547BqKm",
	"GCAkNpk69x9dPo9oapYeWx9ieviUJxvK8HsEzV7pceU6",
	"9awwTAdwtrdhdfLF4GQvLJNnABLaJJ3CsfHFCcDyY4We",
	"6X5KLLfCoCXwFpr5bC1WC7D4zkpgKfTmvPQSwfGgnfTT",
	"CByLkvBRZya16Eu8xkrLsYyZH8ZnTW8htLCcr853ctKS",
	"HRsFSCFntWWLkBeKW18zJqfjpCEPpYhkmjCKPPpgaJPY",
	"2s3P34fEa82ti9kaAV2kvy9qfUoxRtQACBTuWVLLFiE8",
	"828sbDu437wr6nvuMq7YBCGcSX7MKwCNHQyaPLRonhDP",
	"8wXUvU388JosWK5i9zmWJP8dsxLmtsGCUp47nS9Yso8w",
	"GtnuSjoxvNW2gaBxfzq8bevizmtjXN7agrZ1UANSQjVU",
	"8fhCxVtM4ZjLeHS1ZPbR4BxDNSyEB15yEVKEar5jtPpr",
	"HDXcAQ1u4NSPy2ChMdPxAhs8P9ThQUSdmn9HqB1nBhuP",
	"EvirCRx2sAtVUqD2a9V8meZ9ngitTu8PWJyQAS69oq7s",
	"B69SvafJ2QEdCuJV1GCYhj3kPh8pbkQvyWPdPnu3zpWS",
	"AeFeRPZqsAZFVYon3g2fA8XX27ZYBnhrxck2n4DZnhp",
	"HTJEfXkk831BJFZN4AFjccQ221hD4T6XJK8w8GAtParV",
	"4WsWuDPhYhUS22y5Qdd1MBxnCbaFUxpAFMWRpbzq2myi",
	"57EFBcdD2pS2vvJLJQ3akkBcHeiKaVAMQo3HSencBc1f",
	"9eLjvq2tNe9LbC5NJPX4FnxTGZx7rbMJbCQLnNJ38gXu",
	"55HWUJ7gosLDmSSN5j3gShxzEGChPPPH7XMiZG1kwrRV",
	"2gfGMAG49SCetiWX6MwGsoU2yBdeUWCW4sNAUhAvtNA2",
	"EEnvyNZ4skLEdiSzwwB7vWvs2zyGiZovhQdotzvhSe28",
	"94EF39friJjbqh8qg8iXbrKJmxn2na4AFVje3oF21RJL",
	"WktLUp6oFoEuKahctGtYd7fpFyJtJFHUpAg7sKBZTwh",
	"8iQEvMf3JovnDf9qLPAjGYe7VLRWLC25mmQy4cy1om9R",
	"5J5pfy5YGMxP74VRDw2TiQMY2Y8uHLPoaNEe2S5gAtUA",
	"5GRKYKDwhsjARx1NUUUpyEun3ANGtAmRCYXLC4yyVWid",
	"BpAwX5vNLA8QzCCgZHdQSqmj9hW3d7NT16BMKQ7fYZi8",
	"GYBbprmeaoaUCn6QwPS9SntwaCWDrvmxepLusaWUZ645",
	"7gJC2kTpKhx6xEifMdTiPFazKHgK6YS5yg3cdRF6voAi",
	"611DYLZezGmseNXaHLSkff4M1XyCEFyG8KfXLHWabYho",
	"8doqAzcXQzhQ5DkkEvgyP8rRPDksyiYsxQ7XHRZd1j78",
	"2nXcLrurvktTpfyiYLwDPvC2qpjoGogV9tmfEFSSM5Y5",
	"B9qQpxNybcanEs7a7KCMU3bMNzawjatHPxhRrFkie5Gh",
	"3qondpAPfgU2zcuV4uRsEKcvvej4QqDcEyfo8nXUUHy9",
	"7fbJXVk9VDy7NCe8hNY3xQ5u6xoH4ySNtZeoLd1XnygQ",
	"2MRXHrcosRmQdvTqnfo9CTg4J7F1Ayp9cU997TY5Vvjm",
	"3eJkcWY3nje9TYbGfte1n4gRGZrkKrJwVSX1oEUPknNr",
	"8GcekU711bxqLgwtayRebgUPLPQexbUAvoTfr2i4SmS5",
	"8v5o8z7ahTmqz5pBGpvAojCKLyzPzwJPr6oiHgATKcZB",
	"7otqBd9uK2WdLvwJ4eTTbQjwo1Fn9Xi9sfSB1D5qumCp",
	"4D8WFwMx47zN7hkQHvikxvQWiREpdU5Pt2RVoQKHAUif",
	"BchQnzicxp9LCj7zuxFMi4LyjMq1j7NyzYKWqHxoJQhf",
	"HJb8tSRy3Ca3XPoWUTxgGUHD13jVSGFy97HAr8PDytLm",
	"4SxbPp4zCYQPw9KfVB2k86oLZHDdqobGfyjVEiPv8Ebc",
	"3tQU51uNa4pjEdTpZgb2K1wL18cG4dgiyCbmyUpniNZ2",
	"3oMqYtK98qd3PRWvhaUnppESboLZf8gecnYeLxNWULNn",
	"7Qyy6h5CyHEVLQ4bB28r86huSwBzq4fcfNtuyRAQfhWp",
	"EawKJFRTQiBpeuzCk3irFhZJsxVNiUycAobqaLWdJSqM",
	"3jXASE6adPrYi9PKS94brcQ6awdAJcrqqQLSeoUndvXQ",
	"31PivPG33QPMhgQCR24dHu6mAev6o28rsBfjYQHAAvZV",
	"H2gLdqU4axsBH6mcoprCx8hEvMCEFg9bkhDGJH7WygeL",
	"GjBFF5Qo72ZNAK7rjRQcJG2mttynwURdAFjnVRv7CmCt",
	"HyybGnN1es3zANPQjaCyCUkVaca2JvLiKUszQV55v1AK",
	"DqAW3UA6QPANG7t2bUMgpZw9FywECrxXguCw8osYFy5e",
	"HUaV8JpdYizWLkMF9PtF9Lbf1EbYbsZujBZejh2KcsWm",
	"9PU7atrgNR6YJkyx5UnXfUPRb3V3PCZXgW7bhAhrn3Ba",
	"FYCjPg7cPxgk24yDprAUoHuPGv9JsbiPW5uAZGNVAE75",
	"Ff7rwbgvVixFNH1FY2Jq4aXNi1Na1vShDpYoHS5Rcuuf",
	"AFL43dDaH9cDeoaFwBr771NjXbxhYQMWovGgChYraKch",
	"Ajbui7aBhF1x3rLU9dR532ji7UoQixwYMWQRGtFnAEWs",
	"BYB3Espm4wCcj3j2zQtahFb6iLKPzsiELMt7URRrhw7d",
	"61wUqNUVWT8hT6PTzSLHJ1ooQHXZRDNpM4SYzAEWvXrv",
	"4ohrcXfko67P4HMdoh4FukBVEPE37es6F7Mvngkamfa9",
	"EgbyYvFUhB5eF6Y5kRoe1iMtspa3gAAVNHndqbXQYcyS",
	"FcgujJ1K8XbTsGYKrHdyC6wNuxSJHpGiqWLMAuC4HMbi",
	"24sMd9Vfejwu5QTqyu8ig9GbAuEYrbP1yLL6w12Actbg",
	"28jZMYDL4edWKAqY93bhCD3b88bbCfDpVd2eGMrAUdmV",
	"2GcoTj7fehpce8ohnUzira1c36DNUz7kYaUPtU4cBphr",
	"2W6cuQHb14NFZzmhKdgXojRVLRBCfUfVBVrmqThL5RQy",
	"2myNtgxRzLeKCpy87LCzHKVYRg9eZHnsVU8P5hQePjmz",
	"2ttrXT1FXy41gv6D1TL8rM48pYhaaftmZEmH1gDiXRH4",
	"3F7A9b26m5CTxo7GjHCteZzPRcWZ861hhDxra4RQ7QVq",
	"3N2tWvKpak7Geev6trzZz1NjJ5x6j1xPGWotYkN1D4DR",
	"3vLW8KH8RxnfYwVbMGt2K5bDMEvLbfq3gvPvhiVVBf2d",
	"4D8HhpP8tH66vP2co4CDNoptsWKtK9HKBcxfGtxxG3UF",
	"4K8x8sSzipxPwFeYYPSGHS9Ut1Dhb1chpi12khoV8gD3",
	"4QJv95TqLnk5WgtyMHwbVfd5nkWBWNnzHhwnLo64R2yK",
	"4YbVg7KM3zrirMjTWzGWcdUqs8FZ7iFGyKPM73e9RebT",
	"4dMKZcobR8jKUCfFCFhvRvhWp7bbyPkFdr9791RmCUMs",
	"4vFYzx6TcbYkK6YmGun9sX2fgA23HqpCrwB8H87YFovG",
	"5AZq2rLWNryMx3V2bPBJx2TaqERQZ6TYJS5y2SXJacpz",
	"5GVyVAs6nyP52ykqFdgbw1Hg8ujdLHspT3BTmEvtEUds",
	"5ogCSE2JMnro2g6T1VNsoDirHCJ4Cs696yqQUMZULqbm",
	"5y1hmjnA6ZfhFEQvmostj27NgzAzwEYMatMtSCTvhMAo",
	"67hg4ypkiec21B6MTQBpEDgmgiUWxVChJJnnPH9VfvBr",
	"6GoLiLLHHojKX3p6abLcr5tZLWHPjY9Tto7HPMJLKNjd",
	"7MqePDdydYafZPFJHEzCyKozsfF9WcjzpkmcSQrjzMeV",
	"7RdGYyumpRR9P3j95wuxAUMXteofJgf7QMrEyfgCvkhZ",
	"7W3JpG7W7LSfNwaSKSfBjrsiXSEA3NfftCtMbzKfmT5a",
	"7Wjjv7qUrTUtSGAzHiKGJ8owQyvBUsZrBqJ3VXyzRqUg",
	"7XcURKJwFncTEAez2JoF4QvUu4NhwAx8kXbst7xoWH2Y",
	"7dwKdt7PrxYXYUghQikWuVEAXQCQNR2GZkVE6Avy6om7",
	"7mDZPmFVexfzvmwZLDFEfZhjMRZb3fkbS1A7kQJdSGMP",
	"7miMarGfr1Lg4K7p7VeqBo6FcoxDvtkCFAq4DccCB5fe",
	"7vqyTxrXqKdkC24waXgQ5fjkbsmmS6jDFJkwXdLwgmHb",
	"82qeXAbW2uGW9dgkB8ZUB3YskVCdAL3whcnrdEUKz9jC",
	"8491unozSmyxf4ZSu5dimvS64SqQkDEyPYZg2EFatUkB",
	"8KF4fbXf2Uvh1TvZDntZgZRpp54cx2dqt9psxp3MG25b",
	"8WNTMrfwGvXoYvWty6cUu5uwtkaJZK5p4ZMhQbpzE6cf",
	"8gBQE261gQ3n57ZNLRFZBnwioVRGGZSeUyP6d1YTnBoL",
	"8hEJhd2aFYHJ78KLQ9ZVifCuJFNXHcGw7ZqKBqyxYFb2",
	"8xaVYRxo4zHzfiT4GJT1GT87gnRewnPdpVPegwU1XFvT",
	"8xjQnntwvxAcoj1ZpYfUNx1Esw4dAsqXqopcB5MRCqbj",
	"9AjXFcUjBSQi4pzrTUpM95qx5iDRJs4jmnFoC3Q2fzgA",
	"9VCz3j7yqWRSJJ1t2QFq7p2jC5ACz6kyCQSL6RfZthLM",
	"9YqDfGcXYsjjnKjLxM8riqyVBCQ8cXwZQANJ7o8FbayQ",
	"9YvgYrqD2N6gYkvvJsS6JFGytyyFPdVXXgR8HDobHKjn",
	"9ZsVcwcYaB8sis9gP9ThRbegynFEXDd3iZXKhw2qKuyh",
	"9ozhQxVPn2Yp4KH1xE3aHk2SJPaT2ZgXXxQHWZp6wStt",
	"A5W3bry1yB6D6r7ULii6hcNrBwvk1p1eKDhpYkRNjSJt",
	"ANRLMwz2w6UNsqv2XtJwXBPaKqkLx87qfCazNdVqrMqy",
	"ASCg47DcdxqGJFRTpsguiNKAKaAAzvm53NAZWeeoQZHh",
	"Az1C8kC6tj29kGH1R5zFmELQNG1SgeVxH4mhazxsdeEh",
	"BJ4rSMPGWraCbWKRzJ1cWHjTiA6rwcdH9qaqFLTpRDE6",
	"BJEUWTgYZKF84CpqnJUzsVNiE1AcBzoFSXnygrQzbEpr",
	"BZALYiHSyuTCBLRYHtuCHgitKosYdiny9F1URmHJJeM8",
	"BdKxWwkNixHEr3a1rHVkk3nUWmDhZH5wDwjPUt8ypFn8",
	"BojgvBit1itbYtFqjHj8U156gAFuMxzhoi3SdENm7wkf",
	"BsKBTYyyYfVm5cJaHXRZewYF7dcU18aNGF9EbLJ2y2eo",
	"CAj18bAnpzdztLsfr4gyWEp6h9MkJiV4dbv1v9BvVBMS",
	"CFYKi7La1CMvt5AHuXCRSLuDXPmmquG3PiuDjXNW28GC",
	"CQSuv6PCmEdsTyQifLQMxdsUycVq8AiCDajMX3K4LanU",
	"CeWMupuCm4riTaLjXUYTuyM2B9aK4uHzhdn4rw3JNQ1z",
	"CfptrCmZwxcq2fCtYiXvPLAqjwUVKB8Q6X458A58FKvU",
	"ChTccDRYrD5KLQuWcB36M2N9TK6a7X8kGNNpasaJj5Z4",
	"Chs8S271STS1fAK5usAv1cm33TPG5CXPJ9sR9ZihPbek",
	"CtaSd3XXGLL4rSQkdDHqNQ5QRcug5wJcEHy8z79SJYfE",
	"CwUKVVBRpYvkUKHmMy4hbA8xcuuHpLGBdQMrFM3HN1A7",
	"D8fthUX6FEQFSNyjNyzz9Ur3JEdAyxhxHdYPhpbnUfPn",
	"D9FTTazjDdcfq2PKS5pta54qG67drBW6WgRfkQhwanPZ",
	"DcBCRtSUjjXtQNodJT9Rxa8Ht4kJiHSSMpAy5AUb2Cet",
	"Def5ZPMers4gzdUY6ZAPodbZQaTd3qA2se2nHmy4e6f5",
	"DfN8EJsmvDBJKXNh6DBJEKW93tjtnECRLZnfUGMCBDmA",
	"DmJK19vDxU9Wk7sUP8jfPBqtEHqpL5Zm2QCDGRkiUoP1",
	"Dsf9GdLiRu9c8RifX6ez7PNgzyVCdCr46j45d43EWWCy",
	"E1jjUuEVwW5ahc1JcFzMhTjpFDtLrTjDrFrof4xUGMMZ",
	"EBz2KeBDEvn5VeNjEdK7jTkoJcxmQ7xe6tpKdgrqs7Dk",
	"EX4cBrK1RG2VvQZTBPstzKiUEKzmMXTDLjW9TktYRWn7",
	"EaFm5sQT62QbPVod2y31CrQMjjWTZ98aEJVau7cnpv8X",
	"EqNKVHV3Us5SydS7q1Z2q7wX98EcAM1eFrn7akHTP2BN",
	"FCYLQBV9KBgGf2tNTvAAzXBFDifsNWksirAJJy9ydekY",
	"FEDsVtc6wqpdt5NMef2ifWC3ALvgvM7zn6edbg7HCjUq",
	"FPWChzhUc6RyMFGJmjHJJxRVBedzVphPW2U1wXdfAUsh",
	"FVcPXCSGXsRxrAyRxJeMGoruCCCRFvUK2iZ2zg6mNzyK",
	"Fga2EVoafbTXa7Hqgomniu7PkkMAsS3rDqjQ3adW5Wq6",
	"FqZ2oh66gkjq4bZhrE13xAVupgu8htsPb14maceLwtCx",
	"G11VM4ByWp59vtaPTRsSbdW51mtuZERhZRR6eQGQ5Bnd",
	"G9A4Kk1KdErtDuzgLPeZiS9sweaFKL8bKNzqWkkWq61p",
	"GDFt2hbk1TqQSR8EDoGopwZCiGZ6DkLEVFb3Frqrk9LS",
	"GcsXSmP8qZu5nzFnXwJVAqoTCXKFMFVQDYsrD6y9PmRs",
	"Gn6Lqb9Mz7aVLcRxy24Sy3LkUuteGMChZnffLBwwFCVh",
	"GspVkx5kqdzLFme2iMyXewaAbyts9gX3A3KG3F69FUZq",
	"HCk27Njw7WBoUb5VdB2roBHhQWsWKPx1mRcsrswsAJv8",
	"HEd4cn14Ydi6zSCgppMTgDbJWjZCWeSV4k53SXzm2c7W",
	"HNYjknV3REYYtoZEj1u4wBXUZ9RAZViKhpiiRCqSe1Dp",
	"HdWCgBtULzQk8A5wkg7zdLfa9o3zFKGPUZoDKgJ99xJ7",
	"Hm6rz56MCMDw4AwYSne28LC68XP3H1qxKZU62Zz4AsiM",
	"HvtEBos21fFoKJYi81g49bshfZrjPGr3H8AA8zev2V4X",
	"KaZ8ncr6d2v1pmeCttdGQwGL8M5ERHHF7QrTo99c2qQ",
	"LovQMvFwW8XfkMekAPtx12LTk6rWV4HFFBpubhijYDC",
	"TwnGXeK48oDeA8EseiQKDaBReVC73FJWgG9MCpR219R",
	"UEmRps4SXrtiH5zifdvoizCsneYAS6BkxnJ9oh9Pyvu",
	"c5mzVC8RMH424vhYXkt4GFfgZ1wT9GxsuuHiZvxu2QX",
	"ooFGEiW5FGTLb5CTxbDRd6tQKsgfpAF68eipRKQU268",
	"v53q1eoz9dYQ2czTjsspYy2Jkg5xTDRJd8MMyrRLdUF",
	"Egd9CbFpE1d1MWwmRBQkMxx4MVN2pbQPq6EgHtcnmP7i",
	"7a5kPJHaqJDxMG29uuHyYRxmsp2uCA9GyzqPx4FofPEP",
	"F17QNqe7CV88XFwuv3zCDeZFSSLt9Us9xxGTtzmZqvg3",
	"3yFRi8NMzB3n1yoVUyRCL2Pax6a3cPLfYFbCQAKr9vjn",
	"AAGXUsnbHwb1FBNY3hv1vCw7ioQM7h8yqLuyVUkhKHGW",
	"CJQt9kK24wji6HR3jNQqMgYuWAPApicSzEAxabzxji6s",
	"9nVQpUFhygqgPmBhH2rr6bRdQsWkUyvnEuXxfCC9L8SK",
	"9YEUHFFoEDBHsrWBhQWdRsAArkGosbkYYd4rLy1s86NU",
	"6tQz7VTkryswq9PMVus978BbFRbZWYRjZfuTanpCDstb",
	"FDJh7Ne8RCme9qcRPeVMQLRLebdm2ywz6VsSF9JcGb5H",
	"3uL9u35u6h41A6uYxVcgbuA34JwQqhFnwhjGN3gjoEVz",
	"VjS2MehV6JQrmq5F2gGCBVDHsZg1nL9JQ3cdaoBAHrd",
	"C4aMV53ugy7gHg1AAj2eFfa88zaaCDsPBeedo22kco8n",
	"61LkxwBVHumfyuUpnWqiJ1rjuCU1MY2HchUQ4KC1w9sf",
	"9CHgxeBbDBehnNUgoRsrzivffSX2MCPpvfaaF55KnspK",
	"B6QpzmLdGkHiQe199bQ4bhh4T7XRs3BjbnFSmt2sHdXF",
	"J283BjbfU65XCWeLcZbsvjkmSFHWGVMjC7jQStdZbX3m",
	"4igBxeoWRXwKHN8GbBUwpfE3hJvvxdFS3ewJKTANP2zp",
	"pMuFJ4MqDNajMz771oJxHbqZYk1qWiG4J6vLS35Qxj8",
	"6J4qG3NEEbSSikz6YLwGEyAfLf2ei8iAqQ7VScnmm55D",
	"4g34QuaabKAsPYNb6jg8PEisdiBE7DSuQM2Ts8ZYvnHF",
	"v7MjTETtQqETALErnNefyRY5mUuUBKVJNs7FgqbYQdB",
	"Bhkpc8zg4KPBmhj3u4sXQhnXQXEzHxEr8wQfHx2ABuGu",
	"FARcpZoPzTEempwrxJHS8WWKzuqnL6ps6McLUkxmKP3M",
	"DeWwEHaSiYEHUHKqefyGhrh8Ndb5iY5eviAvmXCpkxkB",
	"HJyks3m3ZPT5hTC4Cv642seSxifyREHvgMP5WAmjimBx",
	"HYCAyV1HVxmxvNxsgQwjLYCJWLj2QVa3SSqsEn5EhRrn",
	"wh2Mq3wpwf7F6ei3iuEMxZJ8rHkjFd9Yf1bihSyDn2g",
	"EXvVMwQyQfVT6CfWwzFGetdAvCMS6xuT7B7XdEeEq2eu",
	"Au9xSz6gMtN5JLNCobKcwcSPdBL9mGpfSezAZQZusxdL",
	"Gm2mW3NBdsHPL9WrQ6qcx5ZFhGVMRXiZ8UnffRCupepE",
	"AUVH7TSyUS92aD6no6YiNQptxhBuKvu6ybA1D92xajgP",
	"CkeC9iuGJX9SWdN1J2U11DejLFqRT8q1cvdcuZeZ8R29",
	"BgPpNbWEiYRXxvHmwZNHoYbHFVYygQNWoZVtYrdiF7y1",
	"BKAnqtA6vqSdY6DbSP2gNzaKvuvecYLJEnsZqguAY91j",
	"AsFiJk1yyE3vABEJmFXQNKbKDYuQikyMr2TWnUwPNE3w",
	"7WZdrjScDd7BwHL5vCiqAcs2Qr7t6JfYjLxfVqgg3hwk",
	"CEN5tbRGS3cFL2mgLkoXcJkXz299baBCEb22UsnrS81C",
	"GMyL8Vgksq46dkWCGSY1tLbpNxRPQgfzk8SpGwArmRbb",
	"FCwbwJDCUrEAqNb48RZsTAG7oRMZJHre1BvU51QuYjV7",
	"2rozRCGJDwE79iZ1zFBSofBoSfyUapAh3z8ztJ66DDHz",
	"3qgMSWKQa5xyAmq1zdggLE69D91VwXTeh6d45qid9N9P",
	"G4vJeiiHdDQcqNWxgXca6WCTg5tKgp6b1vWbVSA5wNFR",
	"8G22Mcfo3VRi6LyrsDDThqbESCSM7MdqSVh7BYpGuAwo",
	"847jjUeNpfjfV1ZypHwjvkd5M7rTkjB8c5yPCTnixJZE",
	"F756c7AeegKT32pSiXexEQFZm8b7RC2b95qkoGVnM1rf",
	"CepRC1hLZSNXwtT2Hkk5dfzcvRvZJ5aLVACUPJYH79H4",
	"6eTr7XuVithVrwG67TZ3eDJrveu94T2TXQN5hjDrDn2K",
	"24fwVWTcJjESqEZJ1VQe3syNsox4ZpczsCJA1AxPuCh4",
	"EBVMViDWVPVJsmGmqCc9ZAA51iLoYtafsC52TNeDBvDX",
	"EDguWkMStNsgGxnJ6NMhSVdMEoYVyMQnTUNQtL6WyUvY",
	"EGWfGrDVpnZzCvPNsaPPxbWj2vkxtFmzJup6CmD9Dppo",
	"DQZR9u4QAUTQaG5ajyfqmsSVAYJApWhny3T7CSqeBaFo",
	"4sVYVuL4MAHvTQGMG9LkeBAHDY3fqmMMCu9FvVBcoAod",
	"33445u4ueqyD8MMNfWhdRZKiWkj8KTej5JNFA7gBLirs",
	"4dNCepub5uaHejWSCvjQud9d9X65jABfAspg6ucj45e5",
	"AuMTMc673XDvWQrVerrGm5g3TdY9fYCkoBCSt3N4Zxoc",
	"GHNBqwM6ywfG4mnp8hdymh9ehwNeyorABfekwuhY2BXw",
	"7ttPoYjC2fdTg2ZeTJM3RQCLjoQKbU1JNV1obBPkXQUm",
	"AdccAdowYyMvimrgbgZ6xGouBneN4aZyFD6RNRsH5eLT",
	"CSCQW3zDjLyBKARXXctsJVtKZYAPuWsm9q39fAEU4o3g",
	"8yrLktHsg9V4n1qQMKMgpiQpagzXoWtvWswSf9SCPwVg",
	"9sN5fFizPPV9WhQnqan62wSntMwvPf4JYwSHX5Q8JLUi",
	"7BokPt5VcbSMczZea8wvxPrhFJmq9uM6dDLqiomiWzGZ",
	"4xyCV4LEcZwEPECRHvNrWdnsnv8EJGhNve2t6QZMmagb",
	"AS4uvF14TQro7KtaEEY5b3j6u4xqEJ8RudH3A5gedfpo",
	"52e5Kym9qG5Af9YmQWhUTo4FGSQKxHxySoasTQ8mzz2z",
	"3PM4wiHtsHkBLrRFJXDRDjCxjxacqbqKtNM1HoNTAwGn",
	"3A86qKvSr5rmdvPGWGLqt7FAEwJotPU2tePde7dj8P9t",
	"BVGkY1uoTDoxE7WNoM6TAaAx3EV6xpjKVg3WhJKjLZmF",
	"4WnLg6hLV2RpXBxGQRQFCogiv7cZmg3rVUY9FxGPdNoF",
	"6w4EpWusb3QjdS4rQwUmmKBwNsswLJCG49gPzGPomvYK",
	"59aK1h7jsomG3gf2jAbs3A3V88ik67Uq1UVKqwzdWpWJ",
	"HD42DaXsVYH4BSyjwHvGRjN6M4iR3hFWbtJRHjgXCBHW",
	"41nkUGZcdwZwx6DeghLXQWnYzFZFXmTYC4kYbJejrcpP",
	"Fdv29XQAJbSU477N7jNj6bnt6vv8kDJMyBVgZNdUxBBZ",
	"C2iiEXxR3H8z9H6Z88YU72RSBfXZfpxhyb9XeXYMrXqa",
	"HBZp6NuQMTH7sppnCDh8PcYKBvcsjWTPEZ9cz2bvaGns",
	"BY46jSLARj68EzpMMfZNDGd9uErTAnFVL4ZQPKhDLT3a",
	"HXFwoM1gdv9uFnrjdGp4x1wshvnvu6sk6F2hQDw2FnRw",
	"GztHmwTyYSgG3eH6xwDYS4ZwsmxgUcNmCDntxeZ6pniA",
	"HKz3Sk6xMCWb3zwkGYJbDBSeC7JMM7Zof3uUMBYw91La",
	"4L9Esaozhc9A4cFZSqRNCLKkyJqq5DWHfRsW2QiboHgV",
	"3XCbtYukv4BYNfn5WKyatsX5Vyes7c82T2v4FHu1Jxqg",
	"2aE9uCBmizfxYbk3gQR41v5aZgFEJEo8tVpzDfPCBxfc",
	"64qEppMLoqhW4NAMWFmWdNuexwZ9EazJBptemNfvjMAd",
	"8hzQ3mU3q52aE7W52EVsT9S4h1VcTWGWUVSoyey2RgVE",
	"ExScJ6cYrbf7Y8Ja3dsPpQFM6zrtzB87oYHqUgcPedK3",
	"AeZarQk7hn3pwz6GaN8aYpe13shVA2GUsWm1iJNi3CkY",
	"BEHX6NsnA4VMiMi8nQfbRdHitQFPS4tR8jjYvUwHtofd",
	"E4BruJ6kD9wFZ7FnPwjDVL7ujdrGJMLhNszJWDFHLFHC",
	"3kQ4a2RKr3b7wn88Cf6vDJbU3yf2NMVc7Za2vvqHjNtz",
	"mRSoYCnRDZBEeec4b3azPUQ2n9aQKYty4tyL6s2Fyp5",
	"8KRajZSnYzf6EHaTAxkznVzjyyH1XA8Fh2ZTBM9d9buZ",
	"62WTiRc2nSYJpaLR7vwspEozM6DRPf4b788zXQH7Do8G",
	"BLNZWsTwqA8bbn2AFo6f9iVGepYjHWwVh7d8pLrK7qui",
	"4TwgL2ArEfpq2jWbobRZQtGpLeRbDmwFq5MkjF2EoKXv",
	"PKhFiUsgLqUkg47sv4WsJjn6wpuVwtFo46565iNReWs",
	"8S3dP9WgeQjtijznB9Tet2VBfpig3JJwCBfp7YFFPDRP",
	"5fva3hpGpRUWbYNWGPDwrP5tg4E8YXCQhZ4vtNj6nnQT",
	"4WPg1Y5T6XD6oXsAFnqRJR9tE3h6Yr1WaYXYaKpfWGES",
	"3MnFKupzWKi3DB1621NQVCLQa8ciqLMvfXNSqbTGB3ks",
	"HA35m7iX6Li1TKQDVVnNS8zwqr2qVzhFyufqXKRgSDL7",
	"6ngiuCMtNne3wJNxgiixWdZwkBmBTQPX53Mf47RQVfNv",
	"8uymFKxJb1vWK9XCZ48iiqKYTN92kkdbpXWdyxKAKtyz",
	"AFwsjQHbtNcf1Wbk6LtDzZTTCafzsahY3Eka3FdwvYf4",
	"B9JJFjs4uEaCDwe5q3kxqKQj8AavspcNUuZkJRLEN3Td",
	"Cdbx15czW7qoabDC483bBkK5Le5LuJhoFdtDPvvrsqNb",
	"6j9Noer4fBWBdyb7MZbxyS9rhyAxL4Ty74bzwzfuZwg8",
	"6TeLTCCLbR1zLQM7emwQ3jbnHhKokketaRNv17LxLzey",
	"6pQVraFC63mfmMGc38965uRDQD2gBziJarJ2bPMzEadr",
	"6Cv7s8vExjUoJTXgmptRNtrkvMGEc77ijCwmPdq7vxSB",
	"8hCnGtWmQom79bQPey9fdytcuY5gVrewGHRaJB8qCUM3",
	"7myj8GxV2eqjRviBdMVnuDZR7jNwdnaa3ei3npvHjaE6",
	"Dc2Qt1EutdTWkk7pcFG1RAVJsUbfgG1Y5AonLPEAKgdo",
	"J22G15Hq17egTTAn9F54wTkH6RcehzunSAsjgmD3YhuA",
	"7DnML2B29NGd1DKyp9VrkVJMLhJgwwKqQtxx8zcmbyw1",
	"HF7FvAZLRMrwWjKBXg3vpdWWhB6iWhFRb81VvEGepYUx",
	"CEYcGauHNDkD4FpUG8mJSx7q45MaWYG4dVHCZRZxqQwm",
	"F53cBRaQdRewLJ36M1CrVbbqohVjSUC71pQWmWLpKqu5",
	"BEGo9RonAWKys7RrVu73fmSEEpVJq4zoKLsRdVPsv6YS",
	"FSRtNzejQu4fwWCYCfKtuh9jj7y8JUu42zhiWAboKbkW",
	"5cd6pfZit17wx4Wb8oVyffKyNmKti9vjHe168eJH4KBh",
	"JD4oAKoXJ8VHAV16DJyN5xudm5uZH85ZYhLdC12t9hqT",
	"GXgtoVvsdjCNzcANVBMgLLZBVxt4R8LCnaZUjck5fgtW",
	"EHp3zaWWkkZywa7bEUKAHyx5L5xmPXQepKiyFkw2MuFZ",
	"AdUA6WKHpTcQd5FvgXCxo392n1JS2RuZipcnmAtwYXZi",
	"EmUfjzBAzkDNjqnqNxFATHqXEjhDFabHBkR71e32s2q9",
	"9Phyx2hTefmmSM3qJrwvKhiFvQz5B53PdUGfr2VChV7q",
	"A3bUSSPXifWWf3YfhRM4WL7W9RwskrR4Eb4djkpF1HcT",
	"8dAV8ZpQ9C7FP88Ay3qumoUn2nk15ivff3gwUZcJBunE",
	"DyCmf8ygAhJ5Gc5J4wVpLiSWrX8aoemTFdXogpnPfNuz",
	"hvaNSzsHHCzq61vEz4gGk34b1ZTHE7GfdK9idTkKDQR",
	"2JwiLa1SHpZQksN34tJtY1ucKTGosPvvqeuWsg5nxN1c",
	"2vCEr74697Pw91H2KWespL3BaNEiAXwLfHv1f1Hewtvt",
	"3oDFiC1fuKkoUZWRxreSpkRq8vEXSD71EJ9TriPXKKR8",
	"CZxME7TsrGGRPtQvEnt4NuexyERv4PWsh7BWhGAEZ6Vm",
	"2LsE93mQocnBWokqqwxCMyDiPTqJTeX5YxMszC7hJ4Fz",
	"8EvTLtZXLz2yJtTkv6qeMRoqDkZL7kMDcV5E4HuDLs1z",
	"FHZ8ahKWkQwEDP2zwGtcESDSoKDCT6k39uFT4Um44SH5",
	"8wG7Q6cTq5TfQ5zuSkY4b7VQpuYDsdYjXqS2gNSyS7UN",
	"28ZU77wj89K7finhmfb5YQSgkLAv9zBccm17k2mjGhVg",
	"CPMT9fMoDGph4Q6T8rD1kPvhL8mRpKvUuQcZW92MvQCC",
	"DV98YnGwGjnnuYgqmnAfN52EobdnmbvzDF1fiFxqN4Ah",
	"7ZFPXvH8xeAEcxtGXUHWu93NrJMPsMs5frR8QWBubrmt",
	"GdnanmSGJGqvZbVGXAeaXztAcYHpjH9wq4M4TW9HU9t2",
	"BV52eJ2BmcU6UE265Pcq1giwkYt5ho4u9pYEje1sjt56",
	"2QDcC3HEvdhzxKRoeqwK4vyQq3ts55n6GkGNtXmYGE3E",
	"ZPL1jK3sjTPbZ4QvN9h7sM9GzWia6JwDV8Cs8ZgSAUF",
	"27Abyizeenjhm9ASZzfeCu9ajfMgNGR8nMFWXEdUZhz6",
	"EvsP9sibUJHPxQCyZnrE8hMBHHAUPhMfjZavQ8DtYW2N",
	"6rLrQtpqYTHPZSeYt3ogrnKAfUKC3mwnQ4ER91hRAe4J",
	"4LwHgQ1QMuAveTD3pH29ADdgNuKMETQZZ6Z9t6wLmJ95",
	"HwPUG7E9mEt8Sig5J6EWDNiXbVd82pgZtxMpceJhhxGw",
	"FyN93BcjTJHAopiRXKw9QemNjjeK6nhfbHsEKLpLpSEw",
	"FAFSWPVBFnvNwHYxeC5PqKmQUZErswXvDcxT3N5nPx9B",
	"5rV2xALvQsvCzNQuwLypCWgoNx9SMChQt9wXhJJyUSj9",
	"6JNduzBgAnS8croV1xTrzZRwjyszFn7VxDrSYi1uYHEp",
	"5v9Q6cx4vzvAARm2ZRLwJDY89uju1NuaB3RcTumdCr1s",
	"5sEN198XPpE2FJQXq41K7LF7UYEeSextHfJAj97SRvdL",
	"vBxkWSTiJdjoGfCQfqgERdGV79rpc3zhhtbMehAMU4Y",
	"7SUd2CagBVzkJVEfoNyihaAvyTyTMJjyq8R2mu4CM7tV",
	"AmYsU6Ao97SUUXyst8kRfq8YMJSRsa32bjJeubhVby62",
	"7wyHkJpCHDZm9GXd115NV5qxzBEBR8PaNupZzxJHiLMp",
	"FRMhgZG5jSAfDfzeFZtMo6AKkojuR72myNTeaxXdnLgY",
	"EsZvgHh3FSeJ1P8zncsjEGafTwx3nSBjB1MM9sMdUium",
	"3K7irx54jVmJRTrPi5XWNbMVj3DrgBikDPoBaUobuw7P",
	"5ynUqaLSEGuJTeRaZX8uKacpntmCpuAojcqvTDkwz91U",
	"FBpqtuBrJERENw543qRkx9dk8pUCQvdqPbVuJZ22C7pL",
	"DYS3mepDkhT62A2aJvnPmreUS74Uec34cnPn85AaZq3k",
}

func generateHolderSeed(resourceID int, prefix []byte) string {
	// Преобразование resourceID в байты
	aid := big.NewInt(int64(resourceID)).Bytes()
	// Формируем базу для сидов
	seedBase := append(prefix, aid...)
	// Вычисляем keccak hash и обрезаем до 32 символов
	hash := sha3.NewLegacyKeccak256()
	hash.Write(seedBase)
	keccakHash := hash.Sum(nil)
	return hex.EncodeToString(keccakHash)[:32]
}

func generateHolderAddress(operator string, resourceID int, program string) (solana.PublicKey, error) {
	programKey, err := solana.PublicKeyFromBase58(program)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("invalid program address: %v", err)
	}
	operatorKey, err := solana.PublicKeyFromBase58(operator)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("invalid operator address: %v", err)
	}
	prefix := []byte("holder-")
	seed := generateHolderSeed(resourceID, prefix)
	holderAddress, err := solana.CreateWithSeed(operatorKey, seed, programKey)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("failed to create address with seed: %v", err)
	}
	return holderAddress, nil
}

func fetchAccountData(client *rpc.Client, account solana.PublicKey, dataSize *uint64) ([]byte, error) {
	offset := uint64(0)
	accountInfo, err := client.GetAccountInfoWithOpts(context.TODO(), account, &rpc.GetAccountInfoOpts{
		Encoding: solana.EncodingBase64,
		DataSlice: &rpc.DataSlice{
			Offset: &offset,
			Length: dataSize,
		},
	})
	if err != nil {
		return nil, err
	}
	if accountInfo == nil || accountInfo.Value == nil || accountInfo.Value.Data == nil {
		return nil, fmt.Errorf("account %s not found", account)
	}

	data := accountInfo.Value.Data.GetBinary()
	return data, nil
}

func checkHolderStatus(holders <-chan solana.PublicKey, resChan chan<- struct {
	account solana.PublicKey
	status  string
}, client *rpc.Client, wg *sync.WaitGroup) {
	defer wg.Done()
	dataSize := uint64(8)
	for {
		select {
		case account, ok := <-holders:
			if !ok {
				return
			}
			acc := struct {
				account solana.PublicKey
				status  string
			}{account, "notexist"}
			data, err := fetchAccountData(client, account, &dataSize)
			if err != nil {
				//log.Printf("Failed to fetch data for account %s: %v", account, err)
				resChan <- acc
				continue
			}
			if len(data) > 0 {
				switch data[0] {
				case 32:
					acc.status = "finalized"
				case 52:
					acc.status = "clean"
				default:
					acc.status = "unmatched"
				}
			}
			resChan <- acc
		}
	}
}

func main() {
	var (
		rpcURL    string
		holderNum int
		workerNum int
		program   string
	)

	flag.StringVar(&rpcURL, "rpc", "", "RPC URL for Solana node")
	flag.IntVar(&holderNum, "holders", 1, "Number of holder accounts per operator (1-32)")
	flag.IntVar(&workerNum, "workers", 16, "Number of parallel workers")
	flag.StringVar(&program, "program", "NeonVMyRX5GbCrsAHnUwx1nYYoJAtskU1bWUo6JGNyG", "Program address")
	flag.Parse()

	if rpcURL == "" || program == "" {
		log.Fatal("RPC URL and program address must be specified")
	}

	if holderNum < 1 || holderNum > 32 {
		log.Fatal("holders must be between 1 and 32")
	}

	client := rpc.New(rpcURL)
	var holderAccounts []solana.PublicKey

	for _, operator := range operatorList {
		for i := 1; i <= holderNum; i++ {
			holderAddress, err := generateHolderAddress(operator, i, program)
			if err != nil {
				log.Fatalf("Error generating holder address: %v", err)
			}
			holderAccounts = append(holderAccounts, holderAddress)
		}
	}

	resChan := make(chan struct {
		account solana.PublicKey
		status  string
	}, len(holderAccounts))

	wg := &sync.WaitGroup{}
	wg.Add(workerNum)

	accountsChan := make(chan solana.PublicKey, len(holderAccounts))
	for i := 0; i < len(holderAccounts); i++ {
		accountsChan <- holderAccounts[i]
	}
	close(accountsChan)

	for i := 0; i < workerNum; i++ {
		go checkHolderStatus(accountsChan, resChan, client, wg)
	}
	wg.Wait()
	close(resChan)

	var statuses = make(map[string]int)
	statuses["finalized"] = 0
	statuses["unmatched"] = 0
	statuses["clean"] = 0
	statuses["notexist"] = 0

	for res := range resChan {
		switch res.status {
		case "finalized":
			statuses["finalized"]++
		case "clean":
			statuses["clean"]++
		case "unmatched":
			statuses["unmatched"]++
		case "notexist":
			statuses["notexist"]++
		}
	}

	fmt.Printf("Total holder accounts: %d\n", len(holderAccounts))
	fmt.Printf("Finalized accounts (32): %d\n", statuses["finalized"])
	fmt.Printf("Clean accounts (52): %d\n", statuses["clean"])
	fmt.Printf("Unmatched accounts: %d\n", statuses["unmatched"])
	fmt.Printf("Not exist accounts: %d\n", statuses["notexist"])
}
