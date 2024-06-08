import React, { useState } from 'react';
import { useLocation } from 'react-router-dom';
import './styles/CreateContainer.css'

function useQuery() {
  const { search } = useLocation();
  return React.useMemo(() => new URLSearchParams(search), [search]);
}

const Response = (props) => {
  const query = useQuery();
  if (props.Error !== "") {
    return <p className="response_error">{props.Msg}</p>;
  } else {
    return (
      <p className="response_success">
        Successfully Created the <a href={`/container?ip=${query.get('ip')}&containerID=${props.Msg}`} target="_blank">container!</a>
      </p>
    );
  }
}

const Modal = (props) => {
  const [command, setCommand] = useState('');
  const [response, setResponse] = useState({});
  const query = useQuery();
  const handleCloseModal = () => {
    props.handleCloseModal();
  };

  const handleSubmit = async () => {
    if (command.trim() === '') {
      setResponse('Please enter args.');
    } else {
      let response = await fetch(`http://${process.env.REACT_APP_API_SERVER_URL}/container/create`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ Image: query.get("imageID"), Ip: query.get("ip"), Args: command.trim() })
      });
      let data = await response.json();
      setResponse(data);
    }
  };

  return (
    <div className="modal">
      <div className="modal-content">
        <label htmlFor="command">Args:</label><br />
        <textarea
          id="command"
          name="command"
          rows="10"
          cols="50"
          value={command}
          onChange={(e) => setCommand(e.target.value)}
        /><br />
        {response && <Response Msg={response.Msg} Error={response.Error} />}
        <button className="close-btn" onClick={handleCloseModal}>
           Close
        </button>
        <button className="submit-btn" onClick={handleSubmit}>
          Submit
        </button>
      </div>
    </div>
  );
};

export default Modal;
