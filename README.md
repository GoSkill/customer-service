# customer-service
"Cетевой многопоточный сервис для StatusPage"

Сервис собирает, анализирует и представляет данные о работе систем коммуникации: SMS, MMS, VoiceCall, Email, Billing, Support, Incidents.

Для работы требуется симулятор исходных данных.

Порядок работы:
1. из каталога "simulator" запустить симулятор $ \simulator> go run .
2. запустить сервис $ go run .
3. отредактированные данные смотреть на http://127.0.0.1:8282 или с флагом (пример: $ go run . -addr "http://127.0.0.1:8000")

По умолчанию исходные данные приходят:
"MMS" -> http://127.0.0.1:8383/mms
"Support" -> "http://127.0.0.1:8383/support
"Incidents" -> http://127.0.0.1:8383/accendent

Для пулучения из других источников выполнить запуск с флагом:
пример:
$ go run . -mms "http://127.0.0.1:8383/mms" -support "http://127.0.0.1:8383/support" -accendent "http://127.0.0.1:8383/accendent"