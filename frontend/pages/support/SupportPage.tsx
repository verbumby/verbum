import * as React from 'react'
import { Helmet } from 'react-helmet'
export const SupportPage: React.FC = () => {
	const title = 'Падтрымаць праект'
	return <>
		<Helmet>
			<title>{title}</title>
			<meta name="description" content={title} />
			<meta property="og:title" content={title} />
			<meta property="og:description" content={title} />
			<meta name="robots" content='index, follow' />
		</Helmet>
		<div className='mx-1 mb-3'>
			<h4>{title}</h4>
			<p></p>
			<h5>Шаноўныя сябры!</h5>
			<p></p>
			<p><a href="https://verbum.by/">Verbum.by</a> можна падтрымаць фінансава праз <a href="https://www.buymeacoffee.com/verbum.by" target="_blank" rel="noopener noreferrer"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy Me A Coffee" style={{height: '30px', width: '120px'}} /></a>. Гэта дазволіць пакрываць аплату сервера і даменнага імя, а таксама дадасць матывацыі падтрымліваць і паляпшаць праект.</p>
			<p>Таксама заклікаю вас падтрымаць перакладчыцкую суполку <a href="https://by-reservation.com/" target='_blank' rel="noopener noreferrer" style={{color: 'inherit'}}>«МТГ «Rэ<span style={{color: 'red'}}>З</span>ервацыЯ»</a> праз <a href="https://www.patreon.com/by_reservation" target='_blank' rel="noopener noreferrer">Patreon</a> ці <a href="https://boosty.to/by_reservation" target='_blank' rel="noopener noreferrer">Boosty</a>. Суполка зрабіла і працягвае рабіць агромністы ўнёсак у развіццё гэтага сайта.</p>
			<p>Дзякуй вялікі, што карыстаецеся <a href="https://verbum.by/">verbum.by</a>!</p>
			<p>З павагай,<br/>
			Вадзім.</p>
		</div>
	</>
}
