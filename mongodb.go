
package main

import (
  "github.com/jeffail/gabs"
  "fmt"
  "log"
  "strings"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

type MongoDB struct {
  host string;
  port uint16;
};

type serverStatusResult struct {
  Host            string
	Connections     struct {
		Current      uint32
		Available    uint32
	}
  Ok uint16
}

type replsSetStatusResult struct {
  Set          string
  MyState      uint16 `bson:"myState"`
  Members      []struct {
    Name         string
    Health       uint16
    State        uint16
    StateStr     string `bson:"stateStr"`
    Uptime       uint64
    Self         bool
    PingMs       uint64 `bson:"pingMs"`
  }
}

type currentOpResult struct {
  Inprog      []struct {
    Opid          uint64
    SecsRunning   uint64 `bson:"secs_running"`
    Op            string
    Ns            string
    Client        string
  }
}

func (this MongoDB) Monitor() (string, error) {
  session, err := mgo.Dial(fmt.Sprintf("%s:%d", this.host, this.port))
  if err != nil {
    return "", err
  }

  var serverStatus serverStatusResult
  session.Run("serverStatus", &serverStatus)

  var replSetStatus replsSetStatusResult
  session.DB("admin").Run("replSetGetStatus", &replSetStatus)

  var currentOp currentOpResult
  query := bson.M{}
  coll := session.DB("admin").C("$cmd.sys.inprog")
  coll.Find(query).One(&currentOp)


  atollReport := gabs.New();
  atollReport.SetP("mongodb", "id");
  atollReport.SetP("MongoDB", "name");
  atollReport.ArrayP("report.items");

  // Main state
  state := "ok";
  if len(replSetStatus.Members) > 0 {
    for i, member := range(replSetStatus.Members) {
      atollReport.ArrayAppendP(map[string]interface{}{}, "report.items");
      atollItem := atollReport.S("report").S("items").Index(i);
      if member.Health != 1 {
        state = "error"
      }
      atollItem.SetP(member.Name, "name");
      atollItem.SetP([1]string{"replicaset-member"}, "classes");
      atollItem.SetP(strings.ToLower(member.StateStr), "role");
      atollItem.SetP(member.State, "stats.state.value");
      atollItem.SetP(member.Health, "stats.health.value");
      atollItem.SetP([2]string{"seconds", "duration"}, "stats.uptime.classes");
      atollItem.SetP(member.Uptime, "stats.uptime.value");
      atollItem.SetP([2]string{"miliseconds", "duration"}, "stats.ping.classes");
      atollItem.SetP(member.PingMs, "stats.ping.value");
    }
  }
  atollReport.SetP(state, "report.status.state");

  atollReport.SetP(serverStatus.Connections.Available, "report.stats.availableConnections.value");
  atollReport.SetP(serverStatus.Connections.Current, "report.stats.currentConnections.value");
  atollReport.SetP(len(currentOp.Inprog), "report.stats.numCurrentOps.value");


  var numRunningFor2Seconds uint64 = 0
  var numRunningFor10Seconds uint64 = 0
  var numRunningFor60Seconds uint64 = 0
  var totalRunningSeconds uint64 = 0
  if len(currentOp.Inprog) > 0 {
    for i, _ := range(currentOp.Inprog) {
      op := currentOp.Inprog[i]
      if op.SecsRunning >= 2 {
        numRunningFor2Seconds++;
      }
      if op.SecsRunning >= 10 {
        numRunningFor10Seconds++;
      }
      if op.SecsRunning >= 60 {
        numRunningFor60Seconds++;
      }
      totalRunningSeconds = totalRunningSeconds + op.SecsRunning;
      log.Printf("Op %v, %d", op, op.SecsRunning)
    }
  }

  atollReport.SetP(totalRunningSeconds, "report.stats.currentOpsTotalRunningTime.value");
  atollReport.SetP([2]string{"duration", "seconds"}, "report.stats.currentOpsTotalRunningTime.classes");
  atollReport.SetP(numRunningFor2Seconds, "report.stats.numCurrentOps2Seconds.value");
  atollReport.SetP(numRunningFor10Seconds, "report.stats.numCurrentOps10Seconds.value");
  atollReport.SetP(numRunningFor60Seconds, "report.stats.numCurrentOps60Seconds.value");

  return atollReport.String(), nil;
}
