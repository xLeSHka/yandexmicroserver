const Agent = props => {
	return (
		<div className='container card'>
			<h2 className='expression title-card'>ID: {props.ID}</h2>
			<h3 className='expression title-card'>Status: {props.Status}</h3>
			<p className='expression details'>
				<span>Addres: {props.Address}</span>
				<br />
				<span>Last heartbeat: {props.LastHearBeat}</span>
			</p>
		</div>
	);
};

export default Agent;
