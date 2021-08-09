package websvc

import "fmt"

func (s *server) restore() {
	logError := func(measID string, err error) {
		s.log.info("[restore %s] failed: %v", measID, err)
		s.log.err("[restore %s] failed: %v", measID, err)
	}

	for _, md := range s.mmd {
		// fail early in case there is no bucket
		bck, err := s.database.LookupBucket(fmt.Sprintf(measBucketNameFmt, md.ID))
		if err != nil {
			s.log.info("[restore %s] ignore: bucket not found", md.ID)
			continue
		}

		s.log.info("[restore %s] started", md.ID)

		meas, err := s.mintMeasurement(md.ID, strToBackendIDs(md.BackendIDsStr), md.Description)
		if err != nil {
			logError(md.ID, fmt.Errorf("internal creation failed: %v", err))
			continue
		}

		meas.bucket = bck
		err = s.scheduleWorker(meas)
		if err != nil {
			logError(md.ID, fmt.Errorf("failed to schedule worker: %v", err))
			continue
		}

		// commit successfully restored measurement
		s.log.info("[restore %s] finished, status: %s", meas.ID, meas.Status)
		s.measCache.insert(meas)
	}

	s.log.info("[restore] finished restoring measurement metadata")
}
