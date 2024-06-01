import React from 'react';


export default function NavBar(props) {
  let linksArray = []
  props.linkMap.forEach(data => {
    linksArray.push(
      <li key={data.name} style={{ display: 'inline', marginRight: '20px' }}>
        <a href={data.link} style={{ color: '#fff', textDecoration: 'none' }}>{data.name}</a>
      </li>
    )
  })
  return (
    <nav style={{ backgroundColor: '#333', padding: '10px', marginBottom: '20px' }}>
      <ul style={{ listStyleType: 'none', margin: 0, padding: 0 }}>
        {linksArray}
      </ul>
    </nav>
  )
}