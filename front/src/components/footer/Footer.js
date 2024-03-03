import './style.css';
import gitHub from './gitHub.svg';
import telegram from './telegram.svg';

const Footer = () => {
	return (
		<footer className='footer'>
			<div className='container'>
				<div className='footer__wrapper'>
					<ul className='social'>
						<li className='social__item'>
							<a
								href='https://github.com/xLeSHka/yandexmicroserver'
								target='_blank'
								rel='noopener noreferrer'
							>
								<img src={gitHub} alt='Link to GitHub' />
							</a>
						</li>
						<li className='social__item'>
							<a
								href='http://t.me/rodway'
								target='_blank'
								rel='noopener noreferrer'
							>
								<img src={telegram} alt='Link to' />
							</a>
						</li>
					</ul>
					<div className='copyright'>
						<p>© 2024 All copyrights reserved</p>
					</div>
				</div>
			</div>
		</footer>
	);
};

export default Footer;
