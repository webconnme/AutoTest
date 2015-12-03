package report

import (
    "fmt"
	"log"
    "io"
	"os"
	"path/filepath"
	"time"
)

import (
 	zmq "github.com/alecthomas/gozmq"
)

import (
    "at"
)

func NewAtReportServer( id string, af *at.AtFrame ) (*AtReport, error) {

    ar := &AtReport{}
	
	ar.AF         = af
	ar.ZmqContext = af.ZmqContext
	
	ar.ServerMode = true
	
	ar.cmdPULL ,  _ = ar.ZmqContext.NewSocket(zmq.PULL)	
	ar.cmdPULL.Bind( AR_ZMQ_PROXY_PULL + id )
	
	ar.pollIndex, _ = af.AppendZmqPollItem( zmq.PollItem{ Socket: ar.cmdPULL, Events: zmq.POLLIN} )

	if dir, err := filepath.Abs(filepath.Dir(os.Args[0])); err  != nil {
		return nil, err;
    } else {
	    ar.RootPath =  dir + "/"
	}	

	ar.ReportPath   = "report/"
	ar.FileName  = "rpt"

	ar.rptFile = nil
	
	ar.ResetServer()
	
	return ar, nil
}

func (ar *AtReport) CloseServer() {

   if ar.rptFile != nil {
       ar.rptFile.Close()
   }
   
   ar.cmdPULL.Close()
}

func (ar *AtReport) CallbackServerRxIn( af *at.AtFrame, index int, data []byte )(bool){

	if ar.pollIndex == index {
	    ar.reporter.Printf( string(data) )
	}

    return false	
}

func (ar *AtReport) ReopenReportFile()( error ) {

    if ar.rptFile != nil {
        ar.rptFile.Close()
    }

    t := time.Now()
	ar.FileName  = fmt.Sprintf("rpt-%d-%02d-%02dT%02d:%02d:%02d",
                               t.Year(), t.Month(), t.Day(),  t.Hour(), t.Minute(), t.Second())
	
	rpt_path     := ar.RootPath + ar.ReportPath
    rpt_filename := ar.RootPath + ar.ReportPath + ar.FileName
    fmt.Printf( "\nrpt_file_name = [%s]\n", rpt_filename )	
	
	var err error
	
    if is_been,_ := at.CheckPathExists( rpt_path ); !is_been {
	    if err = os.MkdirAll( rpt_path, 0777 ); err != nil {
	    	 return err
        }
    }
	
    ar.rptFile, err = os.OpenFile( rpt_filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
	    return err
    }

    return nil   
}

func (ar *AtReport) ResetServer()( error ) {

	ar.Deep = 0
	ar.Current = 0

    err := ar.ReopenReportFile()
	if err != nil {
	    return err
	}
	
 	multiLog := io.MultiWriter( ar.rptFile, os.Stdout)
 	
    ar.reporter = log.New( multiLog, ">> ", log.Ldate|log.Ltime )
	
	return err
}

func (ar *AtReport) Record( cmd string , content string )( error ) {

	ar.reporter.Printf( "[%03d][%05d][%-8s] : %s" , ar.Deep, ar.Current, cmd, content )	
	return nil
}

func (ar *AtReport) StartReportOnServer( title string )( error ) {

    if err := ar.ResetServer(); err != nil {
	    return err
	}
	
	ar.Title = title
	ar.Deep = 0
	
	ar.Record( "start" , ar.Title )

    return nil	
}

func (ar *AtReport) EndReportOnServer( )( error ) {

	ar.Record( "end" , "" )
	ar.rptFile.Close()

    return nil	
}


func (ar *AtReport) SetTotalOnServer( value int )( error ) {

	ar.Total = value
	
	content := fmt.Sprintf( "%d", ar.Total )	
	ar.Record( "total" , content )

    return nil	
}

func (ar *AtReport) SetCurrentOnServer( value int )( error ) {

	ar.Current = value

    return nil	
}

func (ar *AtReport) PushDeepOnServer()( error ) {


    return nil	
}

func (ar *AtReport) WriteDocumentOnServer( value string )( error ) {

	ar.Record( "doc" , value )

    return nil	
}

func (ar *AtReport) StartSubOnServer( title string )( error ) {

	ar.Title = title
	ar.Deep++
	ar.Record( "sub" , ar.Title )

    return nil	
}

func (ar *AtReport) EndSubOnServer( )( error ) {

	ar.Record( "subend" , "" )
	ar.Deep--

    return nil	
}

func (ar *AtReport) SetResultPassOnServer( )( error ) {

	ar.Record( "pass" , "" )

    return nil	
}

func (ar *AtReport) SetResultFailOnServer(  reason string )( error ) {

	ar.Record( "fail" , reason )

    return nil	
}

func (ar *AtReport) SetResultErrorOnServer(  reason string )( error ) {

	ar.Record( "error" , reason )

    return nil	
}

