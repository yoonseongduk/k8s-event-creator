# Kubernetes Event Creator

이 프로그램은 Kubernetes 클러스터의 Node와 Service에 대한 테스트 이벤트를 주기적으로 생성하는 도구입니다.

## 사전 준비

1. Go 설치
   ```bash
   # Ubuntu/Debian
   sudo apt-get update
   sudo apt-get install golang-go

   # CentOS/RHEL
   sudo yum install golang

   # macOS
   brew install go

   # Windows 11
   # 1. Go 공식 웹사이트(https://golang.org/dl/)에서 Windows MSI 설치 파일 다운로드
   # 2. 다운로드한 MSI 파일 실행
   # 3. 설치 완료 후 PowerShell이나 명령 프롬프트를 열고 다음 명령어로 설치 확인
   go version

   # VSCode에서 Go 개발 환경 설정
   # 1. VSCode에서 Go 확장 프로그램 설치
   #    - VSCode 확장 마켓플레이스(Ctrl+Shift+X)에서 'Go' 검색
   #    - Microsoft에서 제공하는 Go 확장 프로그램 설치
   # 
   # 2. Go 도구 설치 (VSCode 터미널에서 실행)
   go install -v golang.org/x/tools/gopls@latest
   go install -v github.com/go-delve/delve/cmd/dlv@latest
   go install -v honnef.co/go/tools/cmd/staticcheck@latest
   
   # 3. Windows 환경변수 설정 (필요한 경우)
   # - 시스템 환경 변수에 아래 경로 추가
   # - GOPATH: %USERPROFILE%\go
   # - PATH: %GOPATH%\bin
   ```

2. Kubernetes 클러스터 접근 설정
   ```bash
   # kubeconfig 파일이 ~/.kube/config 위치에 있는지 확인
   ls ~/.kube/config
   ```

## 필요 조건

- Go 1.21 이상
- Kubernetes 클러스터 접근 권한
- `kubectl` 설정 완료

## 설치 방법

1. 저장소 복제
   ```bash
   git clone https://github.com/yourusername/k8s-event-creator.git
   cd k8s-event-creator
   ```

2. 의존성 설치
   ```bash
   go mod download
   ```

3. 빌드
   ```bash
   go build -o k8s-event-creator
   ```

4. 실행
   ```bash 
   ./k8s-event-creator
   ```

## 기능

- Kubernetes 클러스터의 모든 Node에 대한 테스트 이벤트 생성
- Default 네임스페이스의 모든 Service에 대한 테스트 이벤트 생성
- 1시간 간격으로 이벤트 자동 생성
- 로그를 파일(`/var/log/events.log`)과 표준 출력에 동시 기록
