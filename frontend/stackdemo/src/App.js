import { useEffect, useState, useRef } from 'react';

function App() {

	const [greeting, setGreeting] = useState("Greet Someone!");

	const [prevGreetings, setPrevGreetings] = useState([]);

	const greetingName = useRef();
	const greetingLanguage = useRef();


	useEffect(() => {	
		const reqOpt = {
			method: 'get',
			headers: { 'Content-Type': 'application/json' },
		};

		fetch('http://localhost:1337/getGreetings', reqOpt)
			.then(response => response.json())
			.then(data => {
				setPrevGreetings(data)
				console.log(data)
				});
	}, [greeting]);

	function greetSomeone(event) {
		event.preventDefault();

		const reqOpt = {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(
				{ Name: greetingName.current.value, 
				Language: greetingLanguage.current.value }
				)
		};

		fetch('http://localhost:1337/greet', reqOpt)
			.then(response => response.json())
			.then(data => {
				setGreeting(data.Message);
				console.log(data);
				});

	}

	return (
		<div>
			Very cool demo app!!!
			<br/>
			{greeting}
			<form onSubmit={greetSomeone}>
				<select ref={greetingLanguage}>
					<option value="EN">English</option>
					<option value="LA">Latin</option>
					<option value="PL">Polish</option>
					<option value="ES">Spanish</option>
				</select>
				<input ref={greetingName} />
				<input type="submit" value="Greet" />
				<br/>
				Previous Greeting Count: <b>{prevGreetings.Count}</b>
				<br/>
				{(prevGreetings.Greetings != null)?
				prevGreetings.Greetings.map(g => (<li>{g.Message}</li>)):""
				}
				
			</form>
		</div>
	);
}

export default App;

