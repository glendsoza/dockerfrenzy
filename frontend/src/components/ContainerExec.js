import React from 'react';
import Exec from './Exec'
import { useLocation } from 'react-router-dom';

function useQuery() {
    const { search } = useLocation();
    return React.useMemo(() => new URLSearchParams(search), [search]);
}

export default function ContainerExec() {
    const query = useQuery();
    return (
        <Exec url={`ws://${process.env.REACT_APP_API_SERVER_URL}/container/exec?ip=${query.get("ip")}&containerID=${query.get("containerID")}`} readOnly={false} />
    )
}