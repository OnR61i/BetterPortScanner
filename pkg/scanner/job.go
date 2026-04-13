package pkg/scanner

import(
)

type Job struct{
	host net.IP
	port int
	state string
}

func NewJob(host net.IP, port int) Job{
	return Job{
		host: host,
		port: port,
	}
}

func (job *Job) getHost() net.IP{
	return job.host
}

func (job *Job) getPort() int{
	return job.int
}

func (job *Job) getState() (string, error) {
	if job.state != nil{
		return job.state, nil
	}else{
		return "", error.New("state has not been initialized yet!")
	}
}
