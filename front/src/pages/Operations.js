import { sendReqPOST } from '../utils/req.ts';

const Operations = () => {
	const operations = ['plus', 'minus', 'multiply', 'divide'];
	const PATH = 'http://localhost:8082/setOperation';

	const handleSubmit = event => {
		event.preventDefault();
		const form = event.target;
		const formData = new FormData(form);
		let err = 'Pizza';
		operations.forEach(operation => {
			const time = formData.get(operation);
			localStorage.setItem(operation, time);
			let name = '';
			switch (operation) {
				case 'plus':
					name = '+';
					break;
				case 'minus':
					name = '-';
					break;
				case 'multiply':
					name = '*';
					break;
				case 'divide':
					name = '/';
					break;
				default:
					alert(
						'Слушай, я для кого этот код писал со всей логикой? Я тебе сейчас все сломаю >-<'
					);
					throw new Error('НЕ СТОИО СО МНОЙ ШУТИТЬ!');
			}
			const operation_ = {
				Operation: name,
				ExecutionTimeByMilliseconds: String(time),
			};

			if (parseInt(time) <= 0) {
				alert(
					'Слушай, я для кого этот код писал со всей логикой? Я тебе сейчас все сломаю >-<\nЗачем ты вообще в DEVTOOLS полез?'
				);
				throw new Error('НЕ СТОИО СО МНОЙ ШУТИТЬ!');
			}

			sendReqPOST(operation_, PATH).then(response => {
				if (!response.ok && err !== null) {
					err = null;
					alert(
						'Тут вобщем это самое то\nВсе сломалось,\nТвои данные не отправились'
					);
				}
			});
		});
	};

	return (
		<main className='section'>
			<div className='container'>
				<form id='operationsTimeForm' onSubmit={handleSubmit}>
					{operations.map(operation => {
						return (
							<div
								className='operation-container'
								key={operation}
							>
								<label className='title-2'>
									Execution time {operation}
								</label>
								<input
									className='operationInput'
									name={operation}
									type='number'
									placeholder='Enter number'
									required
									min='1'
									defaultValue={localStorage.getItem(
										operation
									)}
								/>
							</div>
						);
					})}
					<input className='btn' type='submit' value='SUBMIT' />
				</form>
			</div>
		</main>
	);
};

export default Operations;
