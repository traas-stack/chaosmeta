import com.sun.tools.attach.VirtualMachine;

import java.io.IOException;

public class ChaosMetaJVMAttacher {
    public static void main(String[] args) {
        String pid = args[0];
        String agentPath = args[1];
        String configPath = args[2];
        VirtualMachine vm = null;
        try {
            vm = VirtualMachine.attach(pid);
        } catch (Exception e) {
            System.out.printf("[error]attach error: %s\n", e.getMessage());
            System.exit(1);
        }

//        System.out.println("attach load before");
        try {
            // Will block until the function logic of agentmain is executed
            vm.loadAgent(agentPath, configPath);
        } catch (Exception e) {
            e.getStackTrace();
            System.out.printf("[error]loadAgent error: %s\n", e.getMessage());
            System.exit(2);
        }
//        System.out.println("attach load after");

        try {
            vm.detach();
        } catch (IOException e) {
            System.out.printf("[error]detach error: %s\n", e.getMessage());
            System.exit(3);
        }
    }
}
