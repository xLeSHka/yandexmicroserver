import React, { useEffect, useState } from 'react';
import { sendReqGET } from '../utils/req.ts';
import ExpressionTag from './Expression.js';
import Error from '../components/error/Error.js';

const Expressions = () => {
	const PATH = 'http://localhost:8082/';
	const [ExpressionArr, setExpressionArr] = useState([]);
	const [newExpression, setNewExpression] = useState('');

	useEffect(() => {
		const fetchData = async () => {
			try {
				const data = await sendReqGET(PATH);
				setExpressionArr(data);
			} catch (error) {
				console.error('Error fetching agents:', error);
			}
		};

		fetchData();
	}, []);

	const handleKeyPress = async event => {
		if (event.key === 'Enter') {
			try {
				// Отправка запроса для добавления нового выражения
				console.log('Добавление нового выражения:', newExpression);

				setNewExpression('');
			} catch (error) {
				console.error('Error adding expression:', error);
			}
		}
	};

	const handleChange = event => {
		setNewExpression(event.target.value);
	};

	return (
		<main className='section'>
			<div className='container'>
				<h2 className='title-2'>Expression</h2>
				<input
					className='expressionInput'
					type='text'
					placeholder='Enter the expression to calculate'
					value={newExpression}
					onChange={handleChange}
					onKeyPress={handleKeyPress}
					required
				/>

				<h2 className='title-2'>Expressions</h2>
				<ul className='projects'>
					{ExpressionArr === undefined ? (
						<Error />
					) : (
						ExpressionArr.map(exp => (
							<ExpressionTag
								key={exp.Id}
								ID={exp.Id}
								Expression={exp.Expression}
								Status={exp.Status}
								CreatedAt={exp.CreatedTime}
								CompletedAt={exp.CompletedTime}
							/>
						))
					)}
				</ul>
			</div>
		</main>
	);
};

export default Expressions;
