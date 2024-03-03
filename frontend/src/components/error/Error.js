import './style.css';
import React, { useState } from 'react';

const Error = () => {
	const fonts = [
		'Arial',
		'Verdana',
		'Tahoma',
		'Courier New',
		'Comic Sans MS',
		'Times New Roman',
		'Georgia',
	];
	const [currentFontIndex, setCurrentFontIndex] = useState(0); // Состояние для хранения индекса текущего шрифта

	const changeFont = () => {
		setCurrentFontIndex(prevIndex => (prevIndex + 1) % fonts.length); // Увеличиваем индекс шрифта на 1, сбрасываем на 0, если достигнут конец массива
	};

	const handleClick = () => {
		alert('Ты в курсе что нажатие не починит ошибку?\n(ʘ ͟ʖ ʘ)'); // Добавляем всплывающее окно при нажатии на кнопку
	};

	return (
		<>
			<h1
				className='error'
				onMouseEnter={changeFont}
				onClick={handleClick}
				style={{ fontFamily: fonts[currentFontIndex] }}
			>
				Error
			</h1>
		</>
	);
};

export default Error;
