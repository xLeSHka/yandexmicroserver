export async function getAllExpressions() {
	try {
		const response = await fetch('http://localhost:8082/expressions');
		if (!response.ok) {
			throw new Error(
				'Failed to fetch expressions. Status code: ' + response.status
			);
		}
		const data = await response.json();
		return data;
	} catch (error) {
		throw new Error('Error fetching expressions: ' + error.message);
	}
}

export async function getAllAgents() {
	try {
		const response = await fetch('http://localhost:8082/agents');
		if (!response.ok) {
			throw new Error(
				'Failed to fetch agents. Status code: ' + response.status
			);
		}
		const data = await response.json();
		return data;
	} catch (error) {
		throw new Error('Error fetching agents: ' + error.message);
	}
}

export async function postOperations(data) {
	try {
		const requestOptions = {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify(data),
		};
		const response = await fetch(
			'http://localhost:8082/setOperation',
			requestOptions
		);
		if (!response.ok) {
			throw new Error(
				'Failed to post operations. Status code: ' + response.status
			);
		}
		const responseData = await response.json();
		console.log(responseData);
	} catch (error) {
		console.error('There was an error!', error);
	}
}

export async function postExpression(expression = '') {
	try {
		const requestOptions = {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({ expression }),
		};
		const response = await fetch(
			'http://localhost:8082/add',
			requestOptions
		);
		if (!response.ok) {
			throw new Error(
				'Failed to post expression. Status code: ' + response.status
			);
		}
		const responseData = await response.json();
		console.log(responseData);
	} catch (error) {
		console.error('There was an error!', error);
	}
}