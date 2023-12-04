import crypto from 'crypto';

const months = [
    { text: 'January', ord: 0 },
    { text: 'February', ord: 1 },
    { text: 'March', ord: 2 },
    { text: 'April', ord: 3 },
    { text: 'May', ord: 4 },
    { text: 'June', ord: 5 },
    { text: 'July', ord:6 },
    { text: 'August', ord: 7},
    { text: 'September', ord: 8 },
    { text: 'October', ord: 9 },
    { text: 'November', ord: 10 },
    { text: 'December', ord: 11 }
];

const days = [
	'Sunday',
	'Monday',
	'Tuesday',
	'Wednesday',
	'Thursday',
	'Friday',
	'Saturday'
];

export function parseICS(file) {
	function expandDates(input) {
		function expandMonth(input) { return months[Number(input) - 1]; }
        const str = `${Number(input.substring(6,8))} ${expandMonth(input.substring(4,6)).text} ${Number(input.substring(0,4))}`;

		return {
			year: Number(input.substring(0,4)),
			month: expandMonth(input.substring(4,6)),
			day: Number(input.substring(6,8)),
			hour: Number(input.substring(9,11)),
			minute: input.substring(11,13),
            jsDate: new Date(Date.parse(str))
		};
	}

	function expandDays(input) {
		switch (input) {
			case 'MO':
				return 'Monday';
			case 'TU':
				return 'Tuesday';
			case 'WE':
				return 'Wednesday';
			case 'TH':
				return 'Thursday';
			case 'FR':
				return 'Friday';
			case 'SA':
				return 'Saturday';
			case 'SU':
				return 'Sunday';
			default:
				throw new Error(`Unexpected day: ${input}`);
		}
	}

	const events = [{}];
	let i = 0;
	let inEvent = false;
	for (const line of file.split('\r\n')) {
		if (!inEvent) {
			inEvent = line === 'BEGIN:VEVENT';
			continue;
		}

		inEvent = !(line === 'END:VEVENT');

		if (inEvent) {
			if (i >= events.length)
				events.push({});

			const lineArr = line.split(':');
			const k = lineArr[0];
			const v = lineArr[1];
			const e = events[i];
            let count;


			switch (k) {
				case 'DTSTART;TZID=America/Chicago':
				case 'DTSTART':
					e.start = expandDates(v);
					break;
				case 'DTEND;TZID=America/Chicago':
				case 'DTEND':
					break;
				case 'RRULE':
					for (const f of v.split(';')) {
						const fArr = f.split('=');
						const k2 = fArr[0];
						const v2 = fArr[1];

						switch (k2) {
							case 'FREQ':
								e.recurring = v2.toLowerCase();
								break;
							case 'BYDAY':
								e.days = v2.split(',').map(expandDays);
								break;
                            case 'COUNT':
                                count = Number(v2);
                                break;
                            case 'UNTIL':
                                e.end = expandDates(v2);
                                break;
                            case 'WKST':
                                break;
							default:
								throw new Error(`Unexpected key in RRULE: ${k2}`);
						}
					}
					break;
				case 'SUMMARY':
					e.title = v;
					break;
				case 'CREATED':
					e.created = expandDates(v);
					break;
				case 'LAST-MODIFIED':
					e.updated = expandDates(v);
					break;
			}

            if (count) {
                const d = new Date(e.start.jsDate.getTime());
                d.setDate(d.getDate() + count*7+1);
                e.end = {
                    year: d.getFullYear(),
                    month: months[d.getMonth()],
                    day: d.getDate(),
                    hour: 0,
                    minute: 0,
                    jsDate: d
                };
            }
		} else {
			if (events[i].recurring == null)
				events[i].recurring = 'never';

			i++;
		}
	}

	return events;
}

// TODO (jacob): There's a lot to consider here
export function generateHTML(events) {
	function htmlEscape(str) {
		return String(str)
			.replace(/&/g, '&amp;')
			.replace(/</g, '&lt;')
			.replace(/>/g, '&gt;')
			.replace(/"/g, '&quot;');
	}
	function upcomingEvents(events) {
		return '';
	}

	function recurringEvents(events) {
		//const months = [...new Set(events.map(e => e.start.month.text))];
		//months.sort((a,b) => b.ord - a.ord);

		let acc = `<table class="schedule">
	<tr>
		<td></td>
		<td>Sunday</td>
		<td>Monday</td>
		<td>Tuesday</td>
		<td>Wednesday</td>
		<td>Thursday</td>
		<td>Friday</td>
		<td>Saturday</td>
	</tr>`;


        const today = (new Date()).setHours(0,0,0,0);
        events = events.filter(el => {
            const end = el?.end?.jsDate?.getTime();
            const start = el?.start?.jsDate?.getTime();

            return end == null 
                || (end > today && start <= today);
        });
		events.sort((a, b) => -(b.start.hour+b.start.minute/60) + (a.start.hour+a.start.minute/60));
		const hours = [...new Set(events.map(e => e.start.hour))];
		const minutes = [...new Set(events.map(e => e.start.minute))];
		for (const hour of hours) {
			const eventsThisHour = events.filter(e => e.start.hour === hour);
			for (const minute of minutes) {
				const eventsThisMinute = eventsThisHour.filter(e => e.start.minute === minute);
				if (eventsThisMinute.length === 0)
					continue;

				let adjustedHour, ap;
				if (hour > 12) {
					adjustedHour = hour - 12;
					ap = 'PM'; 
				} else if (hour === 12) {
					adjustedHour = hour;
					ap = 'PM';
				}else {
					adjustedHour = hour;
					ap = 'AM';
				}

				acc += `\n\t<tr>\n\t\t<td>${adjustedHour}:${minute} ${ap}</td>`;

				for (const day of days) {
					let e = eventsThisMinute.find(e => e.days.includes(day));
					acc += `\n\t\t<!--${day}--><td>${htmlEscape(e?.title ?? '')}</td>`
				}
				acc += '\n\t</tr>';
			}
		}
		acc += '\n</table>';

		return acc;		
	}

	return {
		upcoming: upcomingEvents(events.filter(el => el.recurring === 'never')),
		recurring: recurringEvents(events.filter(el => el.recurring === 'weekly'))
	};
}

export function generateSignIn(classes, students) {
	const today = days[new Date().getDay()];
	
	const todaysClasses = classes.filter(c =>
		c.days.includes(today));

	console.log(`<div id="body">
<h1>Sign In</h1>
<form name="signIn"
	  action="~ROOT/api-local/sign-in" 
	  method="post">
	<div>
		<label for="student">Student</label>
		<select name="student" required>
			<option></option>`);

	for (const s of students)
		console.log(`<option value="${s.id}">${s.first_name} ${s.last_name}</option>`);

	console.log(`</select>
		</div>
		<div>
			<label for="class">Class</label>
			<select name="class" required>
				<option></option>`);

	for (const c of todaysClasses)
		console.log(`<option value="${c.id}">${c.display_title}</option>`);
	
	console.log(`</select>
		</div>
		<input value="Submit" type="submit"/>
	</form>
	<script>window.onunload = () => {}; document.signIn.reset()</script>
</div>`);
}
