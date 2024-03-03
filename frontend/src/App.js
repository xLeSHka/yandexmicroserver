import Navbar from './components/nav/Navbar';
import Footer from './components/footer/Footer';
import './styles/main.css';

import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import ScrollToTop from './utils/scrollToTop';

import Agents from './pages/Agents';
import Operations from './pages/Operations';
import Expressions from './pages/Expressions';
import ExpressionTag from './pages/Expression';

function App() {
	return (
		<div className='App'>
			<Router>
				<ScrollToTop />

				<Navbar />

				<Routes>
					<Route path='/expressions' element={<Expressions />} />
					<Route
						path='/expression/:Id'
						component={<ExpressionTag />}
					/>
					<Route path='/operations' element={<Operations />} />
					<Route path='/agents' element={<Agents />} />
				</Routes>

				<Footer />
			</Router>
		</div>
	);
}

export default App;
