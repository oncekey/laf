import { Module } from '@nestjs/common'
import { ApplicationsController } from './applications.controller'
import { SpecsController } from './specs.controller'
import { CoreModule } from '../core/core.module'
import { ApplicationsService } from './applications.service'
import { PrismaService } from '../prisma.service'

@Module({
  imports: [CoreModule],
  controllers: [ApplicationsController, SpecsController],
  providers: [ApplicationsService, PrismaService],
  exports: [ApplicationsService],
})
export class ApplicationsModule {}
