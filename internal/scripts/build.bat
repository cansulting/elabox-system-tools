@echo off
REM unset go lang env variables
go env -u GOOS
go env -u GOARCH
go env -u GO111MODULE
go env -u CC
go env -u CXX

set PROJ_HOME=..\\..\\..
set ELA_SRC=%PROJ_HOME%\\Elastos.ELA
set EID_SRC=%PROJ_HOME%\\Elastos.ELA.SideChain.EID
set ESC_SRC=%PROJ_HOME%\\Elastos.ELA.SideChain.ESC
set GLIDE_SRC=%PROJ_HOME%\\glide-frontend
set ELA_LANDING=%PROJ_HOME%\\elabox-companion-landing
set ELA_REWARDS=%PROJ_HOME%\\elabox-rewards
set ELA_LOGS=%PROJ_HOME%\\elabox-logs
set ELA_STORE=%PROJ_HOME%\\elabox-dapp-store
set ELA_SETUP=%PROJ_HOME%\\elabox-setup-wizard
set ELA_DASHBOARD=%PROJ_HOME%\\elabox-dashboard
set cos=windows
set carc=amd64
set packageinstaller=packageinstaller
set system_name=system
set packager=packager
set target=%cos%
set arch=%carc%
set gobuild=go build -tags DEBUG
set MODE=DEBUG

REM FLAGS
:parse_args
if "%~1"=="" goto end_parse_args
    if "%~1"=="-o" (
        shift
        set target=%~1
    ) else if "%~1"=="-a" (
        shift
        set arch=%~1
    ) else if "%~1"=="-d" (
        shift
        set MODE=%~1
    )
    shift
goto parse_args
:end_parse_args

echo Optional commandline params -o(target) -a(arch)
echo eg. -o linux -a arm64
echo OS=%target%
echo Arch=%arch%

REM release mode?
echo Build Mode: 1 - RELEASE, 2 - STAGING, 3 - Default - DEBUG (leave empty if DEBUG)
set /p mode=
if "%mode%"=="1" (
    set MODE=RELEASE
    set gobuild=go build -ldflags "-w -s" -tags RELEASE
) else if "%mode%"=="2" (
    set MODE=STAGING
    set gobuild=go build -tags STAGING
)
echo Mode=%MODE%
echo.

REM where binaries will be saved
go env -w CGO_ENABLED=1
echo cgo enabled

REM Questions
echo Rebuild elabox landing page? (y/n)
set /p answerLanding=
echo Rebuild elastos binaries? (y/n)
set /p answerEla=
if exist "%ELA_LOGS%" (
    echo Rebuild logging service? (y/n)
    set /p answerLog=
)
echo Rebuild Glide? (y/n)
set /p answerGlide=
echo Rebuild elastos dapp store? (y/n)
set /p answerDstore=
echo Rebuild Setup Wizard? (y/n)
set /p answerSetup=
echo Rebuild Dashboard? (y/n)
set /p answerDashboard=

if exist "%ELA_REWARDS%" (
    echo Rebuild Rewards? (y/n)
    set /p answerRewards=
)

REM build packager
set buildpath=..\\builds\\%target%
echo Building %packager%
mkdir %buildpath%\\packager
echo %gobuild% -o %buildpath%\\packager\\%packager% ..\\cwd\\%packager%
%gobuild% -o %buildpath%\\packager\\%packager%.exe ..\\cwd\\%packager%

REM build system binaries
if "%target%"=="windows" (
    REM windows intel
    if "%arch%"=="386" (
        go env -w CXX=i686-w64-mingw32-g++ 
        go env -w CC=i686-w64-mingw32-gcc
    ) else (
        REM windows amd
        go env -w CXX=x86_64-w64-mingw32-g++ 
        go env -w CC=x86_64-w64-mingw32-gcc
    )
)
go env -w GOOS=%target% 
go env -w GOARCH=%arch%
echo Building %packageinstaller%
mkdir %buildpath%\\%packageinstaller%\\bin
%gobuild% -o %buildpath%\\%packageinstaller%\\bin\\%packageinstaller%.exe ..\\cwd\\%packageinstaller%
echo Building Elabox System
%gobuild% -o %buildpath%\\%system_name%\\bin\\%system_name%.exe ..\\cwd\\%system_name%
for /f "usebackq tokens=1 delims=" %%i in (`json.bat %buildpath%\\%system_name%\\info.json "program"`) do set programName=%%i
move %buildpath%\\%system_name%\\bin\\%system_name%.exe %buildpath%\\%system_name%\\bin\\%programName%
set programName=%programName:"=%

REM build account manager
echo Building Account Manager
mkdir %buildpath%\\account_manager\\bin
%gobuild% -o %buildpath%\\account_manager\\bin\\account_manager.exe ..\\cwd\\account_manager
for /f "usebackq tokens=1 delims=" %%i in (`json.bat %buildpath%\\account_manager\\info.json "program"`) do set programName=%%i
set programName=%programName:"=%
move %buildpath%\\account_manager\\bin\\account_manager.exe %buildpath%\\account_manager\\bin\\%programName%

REM build notification
echo Building Notification System
mkdir %buildpath%\\notification_center\\bin
%gobuild% -o %buildpath%\\notification_center\\bin ..\\cwd\\notification_center
for /f "usebackq tokens=1 delims=" %%i in (`json.bat %buildpath%\\notification_center\\info.json "program"`) do set programName=%%i
move %buildpath%\\notification_center\\bin\\notification_center.exe %buildpath%\\notification_center\\bin\\%programName%

REM build package manager
echo Building Package Manager
mkdir %buildpath%\\package_manager\\bin
%gobuild% -o %buildpath%\\package_manager\\bin ..\\cwd\\package_manager
for /f "usebackq tokens=1 delims=" %%i in (`json.bat %buildpath%\\package_manager\\info.json "program"`) do set programName=%%i
move %buildpath%\\package_manager\\bin\\package_manager.exe %buildpath%\\package_manager\\bin\\%programName%

REM build reward if exists
if exist "%ELA_REWARDS%" if "%answerRewards%"=="y" (
    pushd %ELA_REWARDS%\\scripts
    call build.bat %target% %arch% %MODE%
    popd
)

REM build app logs
if "%answerLog%"=="y" (
    pushd %ELA_LOGS%\\scripts
    call build.bat %target% %arch% %MODE%
    popd
)

REM unset env variables
go env -u CC
go env -u CXX

REM build system landing page
echo %answerLanding%
if "%answerLanding%"=="y" (
    echo Building Landing Page
    pushd %ELA_LANDING%
    set NODE_OPTIONS=--openssl-legacy-provider
    call npm install
    call npm run build
    popd
    rd /s /q %buildpath%\\system\\www
    mkdir %buildpath%\\system\\www
    xcopy /e %ELA_LANDING%\\build\\* %buildpath%\\system\\www
)

go env -u GO111MODULE


REM build Glide?
if "%answerGlide%"=="y" (
    echo Building Glide...
    pushd %GLIDE_SRC%
    REM npm install
    REM npm run build
    popd
    rd /s /q %buildpath%\\glide\\www
    mkdir %buildpath%\\glide\\www
    xcopy /e %GLIDE_SRC%\\build\\* %buildpath%\\glide\\www
    packager %buildpath%\\glide\\packager.json
)

REM build dapp store?
if "%answerDstore%"=="y" (
    pushd %ELA_STORE%\\scripts
    call build.bat %target% %arch% %MODE%
    popd
)

REM build setup wizard?
if "%answerSetup%"=="y" (
    pushd %ELA_SETUP%\\scripts
    call build.bat %target% %arch% %MODE%
    popd
)

REM build dashboard?
if "%answerDashboard%"=="y" (
    pushd %ELA_DASHBOARD%\\scripts
    call build.bat %target% %arch% %MODE%
    popd
)

REM Packaging
echo Start packaging...
packager %buildpath%\\%packageinstaller%\\packager.json
packager %buildpath%\\account_manager\\packager.json
packager %buildpath%\\notification_center\\packager.json
packager %buildpath%\\package_manager\\packager.json
packager %buildpath%\\system\\packager.json

echo Done.
