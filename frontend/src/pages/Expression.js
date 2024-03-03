import './../styles/main.css';

const ExpressionTag = props => {
	return (
		<div className='container card'>
			<h2 className='expression title-card'>ID: {props.ID}</h2>
			<h3 className='expression title-card'>
				Expression: {props.Expression}
			</h3>
			<p className='expression details'>
				<span>Status: {props.Status}</span>
				<br />
				<span>Crated at: {props.CreatedAt}</span>
				<br />
				<span>Completed at: {props.CompletedAt}</span>
			</p>
		</div>
	);
};

export default ExpressionTag;
