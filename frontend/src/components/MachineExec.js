import React from 'react';
import { useLocation } from 'react-router-dom';
import Exec from './Exec'


function useQuery() {
    const { search } = useLocation();
    return React.useMemo(() => new URLSearchParams(search), [search]);
}

export default function MachineExec() {
    const query = useQuery();
    return (
        <Exec url={`ws://${process.env.REACT_APP_API_SERVER_URL}/machine/exec?ip=${query.get("ip")}`} readOnly={false} />
    )
}