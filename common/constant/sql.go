package constant

const (
	DeleteLock = "DELETE FROM tb_job_lock WHERE id IS NOT NULL;"
	TimeoutJob = "SELECT t2.*,t1.is_enable,t1.deleted_time,t1.retry,t1.author,t1.email FROM( SELECT * FROM tb_job_info) t1 INNER JOIN tb_job_log t2 ON t1.id = t2.job_id WHERE (t2.execute_status=1) AND (( t2.dispatch_time + INTERVAL t1.timeout SECOND ) < CURRENT_TIMESTAMP)"
	RetryJob   = "SELECT t2.*,t1.is_enable as `enable` FROM tb_job_info t1 INNER JOIN tb_job_log  t2 on t1.id=t2.job_id WHERE (t1.retry-t2.retry)>0  AND execute_status IN (-1,3,4) AND t2.processing_status=0;"
)
