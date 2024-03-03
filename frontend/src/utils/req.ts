export interface Expression {
	Id: string;
	Expression: string;
	Status: string;
	CreatedTime: string;
	CompletedTime: string;
}

export interface Operation {
	Operation: string;
	ExecutionTimeByMilliseconds: number;
}

export interface Agent {
	ID: string;
	Address: string;
	Status: string;
	LastHearBeat: string;
}

export async function sendReqPOST(
	data: Expression | Operation | Agent,
	PATH: string
): Promise<string> {
	try {
		const response = await fetch(PATH, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify(data),
		});

		if (response.ok) {
			const responseData = await response.json();
			return responseData;
		} else {
			console.error('HTTP error:', response.status);
			return '';
		}
	} catch (error) {
		console.error('Error:', error);
		return '';
	}
}

export async function sendReqGET(
	PATH: string
): Promise<Expression | Operation | Agent | void> {
	// Создаем экземпляр XMLHttpRequest
	const xhr = new XMLHttpRequest();

	// Открываем соединение
	xhr.open('GET', PATH, true);

	// Устанавливаем обработчик события загрузки данных
	xhr.onload = function () {
		if (xhr.status >= 200 && xhr.status < 300) {
			let responseData;
			switch (PATH.substring(21)) {
				case '/expressions':
					responseData = JSON.parse(xhr.responseText) as Expression[];
					break;
				case '/operations':
					responseData = JSON.parse(xhr.responseText) as Operation[];
					break;
				case '/agents':
					responseData = JSON.parse(xhr.responseText);
					break;
				default:
					return null;
			}

			// Можно использовать полученные данные здесь
			console.log('Данные:', responseData);
			return null;
		} else {
			console.error('Не удалось получить данные:', xhr.statusText);
			return null;
		}
	};

	// Устанавливаем обработчик события ошибки
	xhr.onerror = function () {
		console.error('Ошибка запроса');
	};

	// Отправляем запрос
	xhr.send();
}
