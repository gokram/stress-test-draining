# stress-test-draining


loop:
for timeout := time.After(time.Second * time.Duration(TIMEOUT_SEC)); ; {

     select {
        case <-timeout:
		fmt.Printf("Logging - Stress test completed Timeout expired\n")
		break loop

        default:

	now := time.Now()
	//currentTime := now.Format("20060102150405")

currentTime := now.Format("2006-01-02T15:04:05.999999999Z07:00")



fmt.Printf("MM_TRACE %s %d %d %s %s %s Generated_log\n ",batchid,TIMEOUT_SEC,i,INDEX,application_name,currentTime)

i++





