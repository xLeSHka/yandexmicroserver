import React, { useEffect, useState } from 'react';
import Agent from './Agent';
import { getAllAgents } from '../utils/req.js';
import Error from '../components/error/Error';

const Agents = () => {
	const [AgentsArr, setAgentsArr] = useState([]);

	useEffect(() => {
		const fetchData = async () => {
			try {
				const agents = await getAllAgents();
				setAgentsArr(agents);
			} catch (error) {
				console.error('Error fetching agents:', error);
			}
		};

		fetchData();
	}, []);

	return (
		<main className='section'>
			<div className='container'>
				<h2 className='title-2'>Agents Cards</h2>
				<ul className='projects'>
					{AgentsArr === undefined ? (
						<Error />
					) : (
						AgentsArr.map(ag => (
							<Agent
								key={ag.Id}
								ID={ag.Id}
								Status={ag.Status}
								Address={ag.Address}
								LastHearBeat={ag.LastHearBeat}
							/>
						))
					)}
				</ul>
			</div>
		</main>
	);
};

export default Agents;
